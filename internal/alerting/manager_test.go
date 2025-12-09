package alerting

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
	"health-monitor/internal/storage"
	"health-monitor/pkg/config"
)

func setupTestManager(t *testing.T) (*Manager, *storage.Database, func()) {
	t.Helper()

	log := zerolog.Nop()

	cfg := config.DatabaseConfig{
		Type: "sqlite",
		Path: ":memory:",
	}

	db, err := storage.New(cfg, log)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	targetRepo := storage.NewTargetRepository(db.DB())
	checkResultRepo := storage.NewCheckResultRepository(db.DB())
	incidentRepo := storage.NewIncidentRepository(db.DB())

	manager := NewManager(targetRepo, checkResultRepo, incidentRepo, log)

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}

	return manager, db, cleanup
}

func createTestTarget(ctx context.Context, t *testing.T, db *storage.Database) *domain.Target {
	t.Helper()

	target := &domain.Target{
		ID:       "test-target",
		Name:     "Test Target",
		Type:     domain.TargetTypeHTTP,
		Enabled:  true,
		Interval: 30 * time.Second,
		Timeout:  10 * time.Second,
		Config: map[string]interface{}{
			"url": "https://example.com",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	targetRepo := storage.NewTargetRepository(db.DB())
	if err := targetRepo.Create(ctx, target); err != nil {
		t.Fatalf("Failed to create target: %v", err)
	}

	return target
}

func TestManager_ProcessCheckResult_IncidentCreation(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	result := &domain.CheckResult{
		ID:             1,
		TargetID:       target.ID,
		Status:         domain.CheckStatusFailure,
		StatusCode:     500,
		ResponseTimeMs: 1000,
		Message:        "Internal Server Error",
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := manager.ProcessCheckResult(ctx, result); err != nil {
		t.Fatalf("ProcessCheckResult failed: %v", err)
	}

	incidentRepo := storage.NewIncidentRepository(db.DB())
	incident, err := incidentRepo.GetOngoing(ctx, target.ID)
	if err != nil {
		t.Fatalf("Failed to get ongoing incident: %v", err)
	}

	if incident == nil {
		t.Fatal("Expected incident to be created, got nil")
	}

	if incident.TargetID != target.ID {
		t.Errorf("Expected target_id %s, got %s", target.ID, incident.TargetID)
	}

	if incident.Status != domain.IncidentStatusOngoing {
		t.Errorf("Expected status ongoing, got %s", incident.Status)
	}

	if incident.FailureCount != 1 {
		t.Errorf("Expected failure_count 1, got %d", incident.FailureCount)
	}
}

func TestManager_ProcessCheckResult_IncidentResolution(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	failureResult := &domain.CheckResult{
		ID:             1,
		TargetID:       target.ID,
		Status:         domain.CheckStatusFailure,
		StatusCode:     500,
		ResponseTimeMs: 1000,
		Message:        "Internal Server Error",
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := manager.ProcessCheckResult(ctx, failureResult); err != nil {
		t.Fatalf("ProcessCheckResult (failure) failed: %v", err)
	}

	successResult := &domain.CheckResult{
		ID:             2,
		TargetID:       target.ID,
		Status:         domain.CheckStatusSuccess,
		StatusCode:     200,
		ResponseTimeMs: 150,
		Message:        "OK",
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := manager.ProcessCheckResult(ctx, successResult); err != nil {
		t.Fatalf("ProcessCheckResult (success) failed: %v", err)
	}

	incidentRepo := storage.NewIncidentRepository(db.DB())
	incident, err := incidentRepo.GetOngoing(ctx, target.ID)
	if err != nil {
		t.Fatalf("Failed to get ongoing incident: %v", err)
	}

	if incident != nil {
		t.Errorf("Expected no ongoing incident, but found one with ID %d", incident.ID)
	}
}

func TestManager_ProcessCheckResult_ConsecutiveFailures(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	for i := 1; i <= 5; i++ {
		result := &domain.CheckResult{
			ID:             int64(i),
			TargetID:       target.ID,
			Status:         domain.CheckStatusFailure,
			StatusCode:     500,
			ResponseTimeMs: 1000,
			Message:        "Internal Server Error",
			CheckedAt:      time.Now(),
			Metadata:       make(map[string]interface{}),
		}

		if err := manager.ProcessCheckResult(ctx, result); err != nil {
			t.Fatalf("ProcessCheckResult failed on iteration %d: %v", i, err)
		}
	}

	manager.mu.RLock()
	state := manager.states[target.ID]
	manager.mu.RUnlock()

	if state == nil {
		t.Fatal("Expected state to exist for target")
	}

	if state.consecutiveFailures != 5 {
		t.Errorf("Expected 5 consecutive failures, got %d", state.consecutiveFailures)
	}

	if state.consecutiveSuccesses != 0 {
		t.Errorf("Expected 0 consecutive successes, got %d", state.consecutiveSuccesses)
	}

	incidentRepo := storage.NewIncidentRepository(db.DB())
	incident, err := incidentRepo.GetOngoing(ctx, target.ID)
	if err != nil {
		t.Fatalf("Failed to get ongoing incident: %v", err)
	}

	if incident == nil {
		t.Fatal("Expected incident to be created")
	}

	if incident.FailureCount != 5 {
		t.Errorf("Expected incident failure_count 5, got %d", incident.FailureCount)
	}
}

func TestManager_ProcessCheckResult_StateTracking(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	failureResult := &domain.CheckResult{
		ID:             1,
		TargetID:       target.ID,
		Status:         domain.CheckStatusFailure,
		StatusCode:     500,
		ResponseTimeMs: 1000,
		Message:        "Error",
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := manager.ProcessCheckResult(ctx, failureResult); err != nil {
		t.Fatalf("ProcessCheckResult failed: %v", err)
	}

	successResult := &domain.CheckResult{
		ID:             2,
		TargetID:       target.ID,
		Status:         domain.CheckStatusSuccess,
		StatusCode:     200,
		ResponseTimeMs: 150,
		Message:        "OK",
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := manager.ProcessCheckResult(ctx, successResult); err != nil {
		t.Fatalf("ProcessCheckResult failed: %v", err)
	}

	manager.mu.RLock()
	state := manager.states[target.ID]
	manager.mu.RUnlock()

	if state.consecutiveFailures != 0 {
		t.Errorf("Expected 0 consecutive failures after success, got %d", state.consecutiveFailures)
	}

	if state.consecutiveSuccesses != 1 {
		t.Errorf("Expected 1 consecutive success, got %d", state.consecutiveSuccesses)
	}

	if state.lastStatus != domain.CheckStatusSuccess {
		t.Errorf("Expected last status to be success, got %s", state.lastStatus)
	}
}

type mockNotifier struct {
	alerts []*domain.Alert
}

func (m *mockNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	m.alerts = append(m.alerts, alert)
	return nil
}

func (m *mockNotifier) Type() string {
	return "mock"
}

func (m *mockNotifier) Validate(config map[string]interface{}) error {
	return nil
}

func TestManager_NotifierRegistration(t *testing.T) {
	manager, _, cleanup := setupTestManager(t)
	defer cleanup()

	notifier := &mockNotifier{
		alerts: make([]*domain.Alert, 0),
	}

	manager.RegisterNotifier(notifier)

	retrieved, err := manager.GetNotifier("mock")
	if err != nil {
		t.Fatalf("Failed to get notifier: %v", err)
	}

	if retrieved != notifier {
		t.Error("Retrieved notifier does not match registered notifier")
	}
}

func TestManager_AlertCreation(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	notifier := &mockNotifier{
		alerts: make([]*domain.Alert, 0),
	}
	manager.RegisterNotifier(notifier)

	failureResult := &domain.CheckResult{
		ID:             1,
		TargetID:       target.ID,
		Status:         domain.CheckStatusFailure,
		StatusCode:     500,
		ResponseTimeMs: 1000,
		Message:        "Internal Server Error",
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := manager.ProcessCheckResult(ctx, failureResult); err != nil {
		t.Fatalf("ProcessCheckResult failed: %v", err)
	}

	if len(notifier.alerts) == 0 {
		t.Fatal("Expected at least one alert to be sent")
	}

	alert := notifier.alerts[0]
	if alert.TargetID != target.ID {
		t.Errorf("Expected alert target_id %s, got %s", target.ID, alert.TargetID)
	}

	if alert.Type != domain.AlertTypeDown {
		t.Errorf("Expected alert type DOWN, got %s", alert.Type)
	}
}

func TestManager_SlowResponseAlert(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	notifier := &mockNotifier{
		alerts: make([]*domain.Alert, 0),
	}
	manager.RegisterNotifier(notifier)

	slowResult := &domain.CheckResult{
		ID:             1,
		TargetID:       target.ID,
		Status:         domain.CheckStatusSuccess,
		StatusCode:     200,
		ResponseTimeMs: 6000,
		Message:        "OK",
		CheckedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if err := manager.ProcessCheckResult(ctx, slowResult); err != nil {
		t.Fatalf("ProcessCheckResult failed: %v", err)
	}

	found := false
	for _, alert := range notifier.alerts {
		if alert.Type == domain.AlertTypeSlowResponse {
			found = true
			if alert.Severity != domain.AlertSeverityWarning {
				t.Errorf("Expected warning severity, got %s", alert.Severity)
			}
			break
		}
	}

	if !found {
		t.Error("Expected slow response alert to be triggered")
	}
}

func TestManager_SSLExpiryAlert(t *testing.T) {
	manager, db, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()
	target := createTestTarget(ctx, t, db)

	notifier := &mockNotifier{
		alerts: make([]*domain.Alert, 0),
	}
	manager.RegisterNotifier(notifier)

	result := &domain.CheckResult{
		ID:             1,
		TargetID:       target.ID,
		Status:         domain.CheckStatusSuccess,
		StatusCode:     200,
		ResponseTimeMs: 150,
		Message:        "OK",
		CheckedAt:      time.Now(),
		Metadata: map[string]interface{}{
			"ssl_expires_in_days": 15,
			"ssl_expires_at":      time.Now().Add(15 * 24 * time.Hour),
		},
	}

	if err := manager.ProcessCheckResult(ctx, result); err != nil {
		t.Fatalf("ProcessCheckResult failed: %v", err)
	}

	found := false
	for _, alert := range notifier.alerts {
		if alert.Type == domain.AlertTypeSSLExpiring {
			found = true
			if alert.Severity != domain.AlertSeverityWarning {
				t.Errorf("Expected warning severity, got %s", alert.Severity)
			}
			break
		}
	}

	if !found {
		t.Error("Expected SSL expiry alert to be triggered")
	}
}
