package alerting

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"health-monitor/internal/domain"
	"health-monitor/internal/storage"
)

// recordingNotifier stands in for a voice provider: it records every call
// instead of placing one. No HTTP, no phone, no provider.
type recordingNotifier struct {
	mu      sync.Mutex
	alerts  []*domain.Alert
	ctxErrs []error

	receipt domain.Receipt
	err     error
	before  func(ctx context.Context)
}

func (r *recordingNotifier) Type() string { return domain.NotifierTypeSMSPilot }

func (r *recordingNotifier) Validate(map[string]interface{}) error { return nil }

func (r *recordingNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	_, err := r.NotifyWithReceipt(ctx, alert)
	return err
}

func (r *recordingNotifier) NotifyWithReceipt(ctx context.Context, alert *domain.Alert) (domain.Receipt, error) {
	if r.before != nil {
		r.before(ctx)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.alerts = append(r.alerts, alert)
	r.ctxErrs = append(r.ctxErrs, ctx.Err())

	return r.receipt, r.err
}

func (r *recordingNotifier) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.alerts)
}

// setupDispatchManager wires a manager with the send journal and registers a
// policy-routed recording notifier.
func setupDispatchManager(t *testing.T, config map[string]interface{}) (*Manager, *storage.Database, *recordingNotifier, func()) {
	t.Helper()

	manager, db, cleanup := setupTestManager(t)
	manager.SetSendRepository(storage.NewNotificationSendRepository(db.DB()))

	rec := &recordingNotifier{receipt: domain.Receipt{ExternalID: "server-1", ExternalStatus: "1", Cost: "1.5"}}
	manager.RegisterNotifier(&domain.NotifierConfig{
		ID:     "voice",
		Type:   domain.NotifierTypeSMSPilot,
		Config: config,
	}, rec)

	return manager, db, rec, cleanup
}

func checkResult(targetID string, id int64, status domain.CheckStatus) *domain.CheckResult {
	return &domain.CheckResult{
		ID:             id,
		TargetID:       targetID,
		Status:         status,
		StatusCode:     500,
		ResponseTimeMs: 1000,
		Message:        fmt.Sprintf("check %d: %s", id, status),
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}
}

func process(ctx context.Context, t *testing.T, manager *Manager, result *domain.CheckResult) {
	t.Helper()
	if err := manager.ProcessCheckResult(ctx, result); err != nil {
		t.Fatalf("ProcessCheckResult failed: %v", err)
	}
}

func TestDispatch_NoCallBeforeThreshold(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 3})
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))
	process(ctx, t, manager, checkResult(target.ID, 2, domain.CheckStatusFailure))

	if rec.count() != 0 {
		t.Fatalf("Expected no calls before the threshold, got %d", rec.count())
	}

	process(ctx, t, manager, checkResult(target.ID, 3, domain.CheckStatusFailure))

	if rec.count() != 1 {
		t.Fatalf("Expected exactly 1 call at the threshold, got %d", rec.count())
	}

	sendRepo := storage.NewNotificationSendRepository(db.DB())
	incidentRepo := storage.NewIncidentRepository(db.DB())

	incident, err := incidentRepo.GetOngoing(ctx, target.ID)
	if err != nil || incident == nil {
		t.Fatalf("Expected ongoing incident, got %v (err %v)", incident, err)
	}

	sends, err := sendRepo.ListByIncident(ctx, incident.ID, "voice")
	if err != nil {
		t.Fatalf("Failed to list sends: %v", err)
	}
	if len(sends) != 1 {
		t.Fatalf("Expected 1 send row, got %d", len(sends))
	}
	if sends[0].Status != domain.SendStatusSent {
		t.Errorf("Expected status sent, got %s", sends[0].Status)
	}
	if sends[0].ExternalID != "server-1" {
		t.Errorf("Expected receipt to be stored, got %q", sends[0].ExternalID)
	}
}

func TestDispatch_RecoveryBeforeThresholdCancelsCall(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 3})
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))
	process(ctx, t, manager, checkResult(target.ID, 2, domain.CheckStatusFailure))
	process(ctx, t, manager, checkResult(target.ID, 3, domain.CheckStatusSuccess))
	process(ctx, t, manager, checkResult(target.ID, 4, domain.CheckStatusFailure))

	if rec.count() != 0 {
		t.Fatalf("Expected no calls after recovery reset the counter, got %d", rec.count())
	}
}

// A slow response is a warning, not a failure. It must never place a call,
// even though it does increment the in-memory consecutive counter.
func TestDispatch_WarningIsNotAFailure(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 3})
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	for i := int64(1); i <= 5; i++ {
		process(ctx, t, manager, checkResult(target.ID, i, domain.CheckStatusWarning))
	}

	if rec.count() != 0 {
		t.Fatalf("Expected no calls for warnings, got %d", rec.count())
	}
}

func TestDispatch_OneCallPerIncident(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 2})
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	for i := int64(1); i <= 6; i++ {
		process(ctx, t, manager, checkResult(target.ID, i, domain.CheckStatusFailure))
	}

	if rec.count() != 1 {
		t.Fatalf("Expected exactly 1 call for the whole incident, got %d", rec.count())
	}
}

// A new incident is allowed to call again once the previous one resolved.
func TestDispatch_NewIncidentCallsAgain(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 1})
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))
	process(ctx, t, manager, checkResult(target.ID, 2, domain.CheckStatusSuccess))
	process(ctx, t, manager, checkResult(target.ID, 3, domain.CheckStatusFailure))

	if rec.count() != 2 {
		t.Fatalf("Expected 1 call per incident (2 total), got %d", rec.count())
	}
}

// A restart drops the in-memory state. The journal and the ongoing incident are
// what keep the second process from calling again — and from leaving the old
// incident unresolved forever.
func TestDispatch_RestartDoesNotCallAgain(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 2})
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))
	process(ctx, t, manager, checkResult(target.ID, 2, domain.CheckStatusFailure))

	if rec.count() != 1 {
		t.Fatalf("Expected 1 call before the restart, got %d", rec.count())
	}

	// Second process: same database, fresh in-memory state.
	restarted := NewManager(
		storage.NewTargetRepository(db.DB()),
		storage.NewCheckResultRepository(db.DB()),
		storage.NewIncidentRepository(db.DB()),
		zerolog.Nop(),
	)
	restarted.SetSendRepository(storage.NewNotificationSendRepository(db.DB()))
	restarted.RegisterNotifier(&domain.NotifierConfig{
		ID:     "voice",
		Type:   domain.NotifierTypeSMSPilot,
		Config: map[string]interface{}{"min_failures": 2},
	}, rec)

	process(ctx, t, restarted, checkResult(target.ID, 3, domain.CheckStatusFailure))

	if rec.count() != 1 {
		t.Fatalf("Expected no extra call after restart, got %d", rec.count())
	}

	incidentRepo := storage.NewIncidentRepository(db.DB())

	incidents, err := incidentRepo.ListByTarget(ctx, target.ID, 0, 0)
	if err != nil {
		t.Fatalf("Failed to list incidents: %v", err)
	}
	if len(incidents) != 1 {
		t.Fatalf("Expected the restart to reuse the ongoing incident, got %d incidents", len(incidents))
	}

	process(ctx, t, restarted, checkResult(target.ID, 4, domain.CheckStatusSuccess))

	ongoing, err := incidentRepo.GetOngoing(ctx, target.ID)
	if err != nil {
		t.Fatalf("Failed to get ongoing incident: %v", err)
	}
	if ongoing != nil {
		t.Fatal("Expected the restored incident to resolve on recovery")
	}
}

// An unknown outcome may mean the call went through, so it is terminal.
func TestDispatch_UnknownResultIsNotRetried(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 1})
	defer cleanup()

	rec.err = fmt.Errorf("%w: connection reset", domain.ErrUnknownResult)

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	for i := int64(1); i <= 4; i++ {
		process(ctx, t, manager, checkResult(target.ID, i, domain.CheckStatusFailure))
	}

	if rec.count() != 1 {
		t.Fatalf("Expected exactly 1 attempt, got %d", rec.count())
	}

	incidentRepo := storage.NewIncidentRepository(db.DB())
	incident, err := incidentRepo.GetOngoing(ctx, target.ID)
	if err != nil || incident == nil {
		t.Fatalf("Expected ongoing incident, got %v (err %v)", incident, err)
	}

	sends, err := storage.NewNotificationSendRepository(db.DB()).ListByIncident(ctx, incident.ID, "voice")
	if err != nil {
		t.Fatalf("Failed to list sends: %v", err)
	}
	if len(sends) != 1 {
		t.Fatalf("Expected 1 send row, got %d", len(sends))
	}
	if sends[0].Status != domain.SendStatusUnknown {
		t.Errorf("Expected status unknown, got %s", sends[0].Status)
	}
}

func TestDispatch_ProviderRejectionIsRecorded(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 1})
	defer cleanup()

	rec.err = fmt.Errorf("SMSPILOT rejected the call (232): number is not whitelisted")

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))

	incident, err := storage.NewIncidentRepository(db.DB()).GetOngoing(ctx, target.ID)
	if err != nil || incident == nil {
		t.Fatalf("Expected ongoing incident, got %v (err %v)", incident, err)
	}

	sends, err := storage.NewNotificationSendRepository(db.DB()).ListByIncident(ctx, incident.ID, "voice")
	if err != nil {
		t.Fatalf("Failed to list sends: %v", err)
	}
	if len(sends) != 1 || sends[0].Status != domain.SendStatusFailed {
		t.Fatalf("Expected 1 failed send row, got %+v", sends)
	}
}

// The check context carries the check deadline and is cancelled by the
// scheduler; the call must not inherit that.
func TestDispatch_CallSurvivesCancelledCheckContext(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 1})
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	target := createTestTarget(ctx, t, db)

	// Cancel the check context from inside the provider call.
	rec.before = func(context.Context) { cancel() }

	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))

	if rec.count() != 1 {
		t.Fatalf("Expected 1 call, got %d", rec.count())
	}
	if rec.ctxErrs[0] != nil {
		t.Errorf("Expected a live context inside the call, got %v", rec.ctxErrs[0])
	}

	// The check context is cancelled by now, so assertions use a fresh one.
	assertCtx := context.Background()

	incident, err := storage.NewIncidentRepository(db.DB()).GetOngoing(assertCtx, target.ID)
	if err != nil || incident == nil {
		t.Fatalf("Expected ongoing incident, got %v (err %v)", incident, err)
	}

	sends, err := storage.NewNotificationSendRepository(db.DB()).ListByIncident(assertCtx, incident.ID, "voice")
	if err != nil {
		t.Fatalf("Failed to list sends: %v", err)
	}
	if len(sends) != 1 || sends[0].Status != domain.SendStatusSent {
		t.Fatalf("Expected the result to be recorded despite cancellation, got %+v", sends)
	}
}

// Policy-routed notifiers must never receive the per-check alert fan-out:
// down on the first failure, consecutive failures, slow responses, recovery.
func TestDispatch_PolicyNotifierSkipsAlertFanout(t *testing.T) {
	manager, db, rec, cleanup := setupDispatchManager(t, map[string]interface{}{"min_failures": 99})
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))
	process(ctx, t, manager, checkResult(target.ID, 2, domain.CheckStatusFailure))
	process(ctx, t, manager, checkResult(target.ID, 3, domain.CheckStatusFailure))

	slow := checkResult(target.ID, 4, domain.CheckStatusSuccess)
	slow.ResponseTimeMs = 9000
	process(ctx, t, manager, slow)

	process(ctx, t, manager, checkResult(target.ID, 5, domain.CheckStatusSuccess))

	if rec.count() != 0 {
		t.Fatalf("Expected the policy notifier to receive nothing from the alert path, got %d", rec.count())
	}
}

// A policy that cannot be parsed must not silently downgrade the channel into
// the per-alert fan-out.
func TestDispatch_InvalidPolicyIsNotRegistered(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	manager.SetSendRepository(storage.NewNotificationSendRepository(db.DB()))

	rec := &recordingNotifier{}
	manager.RegisterNotifier(&domain.NotifierConfig{
		ID:     "voice",
		Type:   domain.NotifierTypeSMSPilot,
		Config: map[string]interface{}{"min_failures": 0},
	}, rec)

	if _, err := manager.GetNotifier("voice"); err == nil {
		t.Fatal("Expected the notifier with an invalid policy to stay unregistered")
	}

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)
	process(ctx, t, manager, checkResult(target.ID, 1, domain.CheckStatusFailure))

	if rec.count() != 0 {
		t.Fatalf("Expected no calls from an unregistered notifier, got %d", rec.count())
	}
}
