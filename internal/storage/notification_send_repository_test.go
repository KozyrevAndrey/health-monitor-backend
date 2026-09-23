package storage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"health-monitor/internal/domain"
	"health-monitor/internal/storage"
	"health-monitor/pkg/config"
)

func setupSendRepo(t *testing.T) (*storage.NotificationSendRepository, func()) {
	t.Helper()

	db, err := storage.New(config.DatabaseConfig{Type: "sqlite", Path: ":memory:"}, zerolog.Nop())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	if err := db.AutoMigrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return storage.NewNotificationSendRepository(db.DB()), func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}
}

// The unique index is the last line of defence against a second phone call:
// two concurrent ticks can both decide a call is due.
func TestNotificationSendRepository_ClaimIsExclusive(t *testing.T) {
	repo, cleanup := setupSendRepo(t)
	defer cleanup()

	ctx := context.Background()

	first := &domain.NotificationSend{IncidentID: 1, NotifierID: "voice", Sequence: 1, Status: domain.SendStatusUnknown}
	if err := repo.Claim(ctx, first); err != nil {
		t.Fatalf("First claim failed: %v", err)
	}
	if first.ID == 0 {
		t.Error("Expected the claim to fill in the row ID")
	}

	second := &domain.NotificationSend{IncidentID: 1, NotifierID: "voice", Sequence: 1, Status: domain.SendStatusUnknown}
	if err := repo.Claim(ctx, second); !errors.Is(err, domain.ErrSendAlreadyClaimed) {
		t.Fatalf("Expected ErrSendAlreadyClaimed, got %v", err)
	}

	// A different notifier, a different incident and the next sequence are all
	// separate claims.
	for _, send := range []*domain.NotificationSend{
		{IncidentID: 1, NotifierID: "other", Sequence: 1, Status: domain.SendStatusUnknown},
		{IncidentID: 2, NotifierID: "voice", Sequence: 1, Status: domain.SendStatusUnknown},
		{IncidentID: 1, NotifierID: "voice", Sequence: 2, Status: domain.SendStatusUnknown},
	} {
		if err := repo.Claim(ctx, send); err != nil {
			t.Fatalf("Expected an independent claim to succeed, got %v", err)
		}
	}

	sends, err := repo.ListByIncident(ctx, 1, "voice")
	if err != nil {
		t.Fatalf("Failed to list sends: %v", err)
	}
	if len(sends) != 2 {
		t.Fatalf("Expected 2 sends for incident 1, got %d", len(sends))
	}
	if sends[0].Sequence != 1 || sends[1].Sequence != 2 {
		t.Errorf("Expected sends ordered by sequence, got %d and %d", sends[0].Sequence, sends[1].Sequence)
	}
}

func TestNotificationSendRepository_Complete(t *testing.T) {
	repo, cleanup := setupSendRepo(t)
	defer cleanup()

	ctx := context.Background()

	send := &domain.NotificationSend{IncidentID: 1, NotifierID: "voice", Sequence: 1, Status: domain.SendStatusUnknown}
	if err := repo.Claim(ctx, send); err != nil {
		t.Fatalf("Claim failed: %v", err)
	}

	completedAt := time.Now()
	send.Status = domain.SendStatusSent
	send.ExternalID = "555"
	send.ExternalStatus = "1"
	send.Cost = "3.5"
	send.CompletedAt = &completedAt

	if err := repo.Complete(ctx, send); err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	sends, err := repo.ListByIncident(ctx, 1, "voice")
	if err != nil {
		t.Fatalf("Failed to list sends: %v", err)
	}
	if len(sends) != 1 {
		t.Fatalf("Expected 1 send, got %d", len(sends))
	}
	if sends[0].Status != domain.SendStatusSent || sends[0].ExternalID != "555" || sends[0].Cost != "3.5" {
		t.Errorf("Unexpected stored send: %+v", sends[0])
	}
	if sends[0].CompletedAt == nil {
		t.Error("Expected completed_at to be stored")
	}

	missing := &domain.NotificationSend{ID: 999, Status: domain.SendStatusSent}
	if err := repo.Complete(ctx, missing); err == nil {
		t.Error("Expected completing an unknown row to fail")
	}
}
