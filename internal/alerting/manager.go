package alerting

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

type Manager struct {
	targetRepo      domain.TargetRepository
	checkResultRepo domain.CheckResultRepository
	incidentRepo    domain.IncidentRepository
	log             zerolog.Logger

	mu        sync.RWMutex
	notifiers map[string]domain.Notifier
	states    map[string]*targetState
}

type targetState struct {
	targetID            string
	consecutiveFailures int
	consecutiveSuccesses int
	lastStatus          domain.CheckStatus
	lastCheckTime       time.Time
	currentIncidentID   *int64
	lastAlertTime       map[domain.AlertType]time.Time
}

func NewManager(
	targetRepo domain.TargetRepository,
	checkResultRepo domain.CheckResultRepository,
	incidentRepo domain.IncidentRepository,
	log zerolog.Logger,
) *Manager {
	return &Manager{
		targetRepo:      targetRepo,
		checkResultRepo: checkResultRepo,
		incidentRepo:    incidentRepo,
		log:             log,
		notifiers:       make(map[string]domain.Notifier),
		states:          make(map[string]*targetState),
	}
}

func (m *Manager) RegisterNotifier(notifier domain.Notifier) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifiers[notifier.Type()] = notifier
	m.log.Info().Str("type", notifier.Type()).Msg("Notifier registered")
}

func (m *Manager) GetNotifier(notifierType string) (domain.Notifier, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	notifier, exists := m.notifiers[notifierType]
	if !exists {
		return nil, fmt.Errorf("notifier type %s not found", notifierType)
	}
	return notifier, nil
}

func (m *Manager) ProcessCheckResult(ctx context.Context, result *domain.CheckResult) error {
	m.mu.Lock()
	state, exists := m.states[result.TargetID]
	if !exists {
		state = &targetState{
			targetID:      result.TargetID,
			lastAlertTime: make(map[domain.AlertType]time.Time),
		}
		m.states[result.TargetID] = state
	}
	m.mu.Unlock()

	target, err := m.targetRepo.Get(ctx, result.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get target: %w", err)
	}

	previousStatus := state.lastStatus

	m.updateState(state, result)

	if err := m.evaluateAlertRules(ctx, target, result, state, previousStatus); err != nil {
		m.log.Error().
			Err(err).
			Str("target_id", result.TargetID).
			Msg("Failed to evaluate alert rules")
	}

	if err := m.manageIncident(ctx, target, result, state); err != nil {
		m.log.Error().
			Err(err).
			Str("target_id", result.TargetID).
			Msg("Failed to manage incident")
	}

	return nil
}

func (m *Manager) updateState(state *targetState, result *domain.CheckResult) {
	state.lastCheckTime = result.CheckedAt

	if result.Status == domain.CheckStatusSuccess {
		state.consecutiveSuccesses++
		state.consecutiveFailures = 0
	} else {
		state.consecutiveFailures++
		state.consecutiveSuccesses = 0
	}

	state.lastStatus = result.Status
}

func (m *Manager) evaluateAlertRules(
	ctx context.Context,
	target *domain.Target,
	result *domain.CheckResult,
	state *targetState,
	previousStatus domain.CheckStatus,
) error {
	rules := m.getDefaultAlertRules(target)

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		if shouldTrigger, alertType := m.shouldTriggerAlert(rule, result, state, previousStatus); shouldTrigger {
			if m.shouldSendAlert(state, alertType) {
				alert := m.buildAlert(target, result, alertType, rule.Severity)
				if err := m.CreateAlert(ctx, alert); err != nil {
					m.log.Error().
						Err(err).
						Str("target_id", target.ID).
						Str("alert_type", string(alertType)).
						Msg("Failed to create alert")
				} else {
					state.lastAlertTime[alertType] = time.Now()
				}
			}
		}
	}

	return nil
}

func (m *Manager) shouldTriggerAlert(
	rule *domain.AlertRule,
	result *domain.CheckResult,
	state *targetState,
	previousStatus domain.CheckStatus,
) (bool, domain.AlertType) {
	if rule.ConsecutiveFailures > 0 && state.consecutiveFailures >= rule.ConsecutiveFailures {
		return true, domain.AlertTypeConsecutiveFail
	}

	if rule.ResponseTimeMs > 0 && result.ResponseTimeMs > int64(rule.ResponseTimeMs) {
		return true, domain.AlertTypeSlowResponse
	}

	if rule.SSLExpiryDays > 0 {
		if expiryDays, ok := result.Metadata["ssl_expires_in_days"].(int); ok {
			if expiryDays <= rule.SSLExpiryDays && expiryDays > 0 {
				return true, domain.AlertTypeSSLExpiring
			}
		}
	}

	if result.Status == domain.CheckStatusFailure &&
		(previousStatus == domain.CheckStatusSuccess || previousStatus == "") {
		return true, domain.AlertTypeDown
	}

	if result.Status == domain.CheckStatusSuccess &&
		previousStatus == domain.CheckStatusFailure {
		return true, domain.AlertTypeUp
	}

	return false, ""
}

func (m *Manager) shouldSendAlert(state *targetState, alertType domain.AlertType) bool {
	lastSent, exists := state.lastAlertTime[alertType]
	if !exists {
		return true
	}

	cooldownPeriod := 5 * time.Minute
	if alertType == domain.AlertTypeSlowResponse {
		cooldownPeriod = 15 * time.Minute
	}

	return time.Since(lastSent) > cooldownPeriod
}

func (m *Manager) buildAlert(
	target *domain.Target,
	result *domain.CheckResult,
	alertType domain.AlertType,
	severity domain.AlertSeverity,
) *domain.Alert {
	alert := &domain.Alert{
		ID:         uuid.New().String(),
		TargetID:   target.ID,
		TargetName: target.Name,
		Type:       alertType,
		Severity:   severity,
		CreatedAt:  time.Now(),
		Metadata:   make(map[string]interface{}),
	}

	switch alertType {
	case domain.AlertTypeDown:
		alert.Message = fmt.Sprintf("Target %s is DOWN", target.Name)
		alert.Description = result.Message
	case domain.AlertTypeUp:
		alert.Message = fmt.Sprintf("Target %s is UP", target.Name)
		alert.Description = "Service has recovered"
	case domain.AlertTypeSlowResponse:
		alert.Message = fmt.Sprintf("Target %s has slow response time", target.Name)
		alert.Description = fmt.Sprintf("Response time: %dms", result.ResponseTimeMs)
		alert.Metadata["response_time_ms"] = result.ResponseTimeMs
	case domain.AlertTypeSSLExpiring:
		expiryDays := result.Metadata["ssl_expires_in_days"]
		alert.Message = fmt.Sprintf("SSL certificate for %s is expiring soon", target.Name)
		alert.Description = fmt.Sprintf("Certificate expires in %v days", expiryDays)
		alert.Metadata["ssl_expires_in_days"] = expiryDays
		alert.Metadata["ssl_expires_at"] = result.Metadata["ssl_expires_at"]
	case domain.AlertTypeConsecutiveFail:
		alert.Message = fmt.Sprintf("Target %s has consecutive failures", target.Name)
		alert.Description = result.Message
	}

	alert.Metadata["check_result_id"] = result.ID
	alert.Metadata["target_type"] = string(target.Type)

	return alert
}

func (m *Manager) CreateAlert(ctx context.Context, alert *domain.Alert) error {
	m.log.Info().
		Str("target_id", alert.TargetID).
		Str("target_name", alert.TargetName).
		Str("type", string(alert.Type)).
		Str("severity", string(alert.Severity)).
		Msg("Alert created")

	m.mu.RLock()
	notifiers := make([]domain.Notifier, 0, len(m.notifiers))
	for _, notifier := range m.notifiers {
		notifiers = append(notifiers, notifier)
	}
	m.mu.RUnlock()

	if len(notifiers) == 0 {
		m.log.Debug().Msg("No notifiers registered, alert not sent")
		return nil
	}

	for _, notifier := range notifiers {
		if err := notifier.Notify(ctx, alert); err != nil {
			m.log.Error().
				Err(err).
				Str("notifier_type", notifier.Type()).
				Str("alert_id", alert.ID).
				Msg("Failed to send notification")
		}
	}

	return nil
}

func (m *Manager) manageIncident(
	ctx context.Context,
	target *domain.Target,
	result *domain.CheckResult,
	state *targetState,
) error {
	if result.Status == domain.CheckStatusFailure {
		return m.handleFailure(ctx, target, result, state)
	}

	if result.Status == domain.CheckStatusSuccess && state.currentIncidentID != nil {
		return m.resolveIncident(ctx, *state.currentIncidentID, state)
	}

	return nil
}

func (m *Manager) handleFailure(
	ctx context.Context,
	target *domain.Target,
	result *domain.CheckResult,
	state *targetState,
) error {
	if state.currentIncidentID != nil {
		incident, err := m.incidentRepo.Get(ctx, *state.currentIncidentID)
		if err != nil {
			return fmt.Errorf("failed to get incident: %w", err)
		}

		incident.FailureCount++
		incident.LastError = result.Message
		incident.LastCheckResult = result

		if err := m.incidentRepo.Update(ctx, incident); err != nil {
			return fmt.Errorf("failed to update incident: %w", err)
		}

		return nil
	}

	incident := &domain.Incident{
		TargetID:         target.ID,
		TargetName:       target.Name,
		Status:           domain.IncidentStatusOngoing,
		StartedAt:        result.CheckedAt,
		FailureCount:     1,
		LastError:        result.Message,
		Severity:         domain.AlertSeverityCritical,
		FirstCheckResult: result,
		LastCheckResult:  result,
	}

	if err := m.incidentRepo.Create(ctx, incident); err != nil {
		return fmt.Errorf("failed to create incident: %w", err)
	}

	state.currentIncidentID = &incident.ID

	m.log.Info().
		Str("target_id", target.ID).
		Int64("incident_id", incident.ID).
		Msg("Incident created")

	return nil
}

func (m *Manager) resolveIncident(ctx context.Context, incidentID int64, state *targetState) error {
	if err := m.incidentRepo.Resolve(ctx, incidentID); err != nil {
		return fmt.Errorf("failed to resolve incident: %w", err)
	}

	state.currentIncidentID = nil

	m.log.Info().
		Int64("incident_id", incidentID).
		Msg("Incident resolved")

	return nil
}

func (m *Manager) getDefaultAlertRules(target *domain.Target) []*domain.AlertRule {
	return []*domain.AlertRule{
		{
			ID:                  "consecutive-failures",
			Name:                "Consecutive Failures",
			Enabled:             true,
			ConsecutiveFailures: 3,
			Severity:            domain.AlertSeverityCritical,
		},
		{
			ID:             "slow-response",
			Name:           "Slow Response Time",
			Enabled:        true,
			ResponseTimeMs: 5000,
			Severity:       domain.AlertSeverityWarning,
		},
		{
			ID:            "ssl-expiring",
			Name:          "SSL Certificate Expiring",
			Enabled:       true,
			SSLExpiryDays: 30,
			Severity:      domain.AlertSeverityWarning,
		},
	}
}
