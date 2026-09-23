package alerting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"health-monitor/internal/domain"
)

const (
	// notificationDispatchTimeout bounds the whole dispatch step once it is
	// detached from the check context.
	notificationDispatchTimeout = 60 * time.Second

	// notificationSendTimeout bounds a single provider call. It is deliberately
	// generous: a voice provider is slower than a chat API.
	notificationSendTimeout = 30 * time.Second

	// notificationCompleteTimeout bounds the follow-up write that records the
	// outcome. Losing this write would leave the send row claimed but unfinished.
	notificationCompleteTimeout = 5 * time.Second
)

// dispatchIncidentNotifications fires policy-routed notifiers (voice calls) for
// the target's ongoing incident.
//
// It runs after manageIncident, on the incident rather than on the alert, for
// three reasons: the alert path notifies on the first failure and cannot honour
// a configurable threshold, its in-memory cooldown would shift the call by
// minutes, and the incident's FailureCount counts only real failures (a slow
// response must not place a call) and survives a restart.
func (m *Manager) dispatchIncidentNotifications(ctx context.Context, target *domain.Target, state *targetState) {
	if state.currentIncidentID == nil {
		return
	}

	entries := m.policyNotifiers()
	if len(entries) == 0 {
		return
	}

	// Detach from the check context before touching the database or the
	// provider: it carries the check deadline (target timeout + 5s), which a
	// timed-out check has mostly spent, and the alert fan-out ran on it first.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), notificationDispatchTimeout)
	defer cancel()

	if m.sendRepo == nil {
		m.log.Error().
			Str("target_id", target.ID).
			Msg("Notification send repository is not configured, policy notifiers disabled")
		return
	}

	incident, err := m.incidentRepo.Get(ctx, *state.currentIncidentID)
	if err != nil {
		m.log.Error().
			Err(err).
			Str("target_id", target.ID).
			Int64("incident_id", *state.currentIncidentID).
			Msg("Failed to load incident for notification dispatch")
		return
	}

	for _, entry := range entries {
		if err := m.dispatchOne(ctx, target, incident, entry); err != nil {
			m.log.Error().
				Err(err).
				Str("notifier_id", entry.id).
				Int64("incident_id", incident.ID).
				Msg("Failed to dispatch incident notification")
		}
	}
}

// policyNotifiers returns a snapshot of the registered policy-routed notifiers.
func (m *Manager) policyNotifiers() []*registeredNotifier {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries := make([]*registeredNotifier, 0, len(m.notifiers))
	for _, entry := range m.notifiers {
		if entry.policy != nil {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (m *Manager) dispatchOne(
	ctx context.Context,
	target *domain.Target,
	incident *domain.Incident,
	entry *registeredNotifier,
) error {
	policy := *entry.policy

	if !policy.AppliesTo(target.ID) {
		return nil
	}

	if incident.FailureCount < policy.MinFailures {
		return nil
	}

	sends, err := m.sendRepo.ListByIncident(ctx, incident.ID, entry.id)
	if err != nil {
		return fmt.Errorf("failed to list previous sends: %w", err)
	}

	if len(sends) >= policy.MaxSends {
		return nil
	}

	if len(sends) > 0 {
		// Repeats are opt-in: an empty interval means "notify once per incident".
		if policy.RepeatInterval == 0 {
			return nil
		}
		last := sends[len(sends)-1]
		if time.Since(last.CreatedAt) < policy.RepeatInterval {
			return nil
		}
	}

	send := &domain.NotificationSend{
		IncidentID: incident.ID,
		NotifierID: entry.id,
		Sequence:   len(sends) + 1,
		// Claimed rows start as unknown, so a crash between the claim and the
		// provider response leaves a terminal row that is never retried.
		Status: domain.SendStatusUnknown,
	}

	if err := m.sendRepo.Claim(ctx, send); err != nil {
		if errors.Is(err, domain.ErrSendAlreadyClaimed) {
			return nil
		}
		return fmt.Errorf("failed to claim notification send: %w", err)
	}

	alert := buildIncidentAlert(target, incident)

	// ctx is already detached from the check deadline by the caller.
	sendCtx, cancel := context.WithTimeout(ctx, notificationSendTimeout)
	defer cancel()

	receipt, sendErr := notifyWithReceipt(sendCtx, entry.notifier, alert)

	completedAt := time.Now()
	send.CompletedAt = &completedAt
	send.ExternalID = receipt.ExternalID
	send.ExternalStatus = receipt.ExternalStatus
	send.Cost = receipt.Cost

	switch {
	case sendErr == nil:
		send.Status = domain.SendStatusSent
		m.log.Info().
			Str("notifier_id", entry.id).
			Str("target_id", target.ID).
			Int64("incident_id", incident.ID).
			Int("sequence", send.Sequence).
			Str("external_id", receipt.ExternalID).
			Msg("Incident notification sent")

	case errors.Is(sendErr, domain.ErrUnknownResult):
		send.Status = domain.SendStatusUnknown
		send.Error = sendErr.Error()
		// Terminal on purpose: the notification may have gone out, and a retry
		// could mean a second phone call.
		m.log.Error().
			Err(sendErr).
			Str("notifier_id", entry.id).
			Int64("incident_id", incident.ID).
			Msg("Incident notification result unknown, not retrying")

	default:
		send.Status = domain.SendStatusFailed
		send.Error = sendErr.Error()
		m.log.Error().
			Err(sendErr).
			Str("notifier_id", entry.id).
			Int64("incident_id", incident.ID).
			Msg("Incident notification rejected by provider")
	}

	completeCtx, completeCancel := context.WithTimeout(ctx, notificationCompleteTimeout)
	defer completeCancel()

	if err := m.sendRepo.Complete(completeCtx, send); err != nil {
		return fmt.Errorf("failed to record notification result: %w", err)
	}

	return nil
}

// notifyWithReceipt prefers the receipt-aware call so provider identifiers can
// be stored, and falls back to plain Notify for notifiers without one.
func notifyWithReceipt(ctx context.Context, n domain.Notifier, alert *domain.Alert) (domain.Receipt, error) {
	if receiptNotifier, ok := n.(domain.ReceiptNotifier); ok {
		return receiptNotifier.NotifyWithReceipt(ctx, alert)
	}
	return domain.Receipt{}, n.Notify(ctx, alert)
}

// buildIncidentAlert renders the alert passed to policy-routed notifiers. It
// describes the incident, not a single check, since the incident is what the
// on-call person is being woken up for.
func buildIncidentAlert(target *domain.Target, incident *domain.Incident) *domain.Alert {
	return &domain.Alert{
		ID:          uuid.New().String(),
		TargetID:    target.ID,
		TargetName:  target.Name,
		Type:        domain.AlertTypeDown,
		Severity:    incident.Severity,
		Message:     fmt.Sprintf("Target %s is DOWN", target.Name),
		Description: incident.LastError,
		CreatedAt:   time.Now(),
		Metadata: map[string]interface{}{
			"incident_id":   incident.ID,
			"failure_count": incident.FailureCount,
			"started_at":    incident.StartedAt,
			"target_type":   string(target.Type),
		},
	}
}
