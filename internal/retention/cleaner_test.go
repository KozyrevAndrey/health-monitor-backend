package retention

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"health-monitor/internal/domain"
	"health-monitor/pkg/config"
)

// fakeResultRepo embeds the interface so only DeleteOlderThan needs overriding.
type fakeResultRepo struct {
	domain.CheckResultRepository
	calls  int
	before time.Time
}

func (f *fakeResultRepo) DeleteOlderThan(_ context.Context, before time.Time) (int64, error) {
	f.calls++
	f.before = before
	return 5, nil
}

type fakeIncidentRepo struct {
	domain.IncidentRepository
	calls  int
	before time.Time
}

func (f *fakeIncidentRepo) DeleteResolvedOlderThan(_ context.Context, before time.Time) (int64, error) {
	f.calls++
	f.before = before
	return 2, nil
}

func TestCleaner_Cleanup_BothWindows(t *testing.T) {
	results := &fakeResultRepo{}
	incidents := &fakeIncidentRepo{}
	cfg := config.RetentionConfig{
		CheckResults:    24 * time.Hour,
		Incidents:       48 * time.Hour,
		CleanupInterval: time.Hour,
	}
	c := New(results, incidents, cfg, zerolog.Nop())

	now := time.Now()
	c.cleanup(context.Background())

	if results.calls != 1 || incidents.calls != 1 {
		t.Fatalf("Expected both deletes called once, got results=%d incidents=%d", results.calls, incidents.calls)
	}
	// Cutoff should be roughly now - ttl (allow a small skew).
	if skew := now.Add(-cfg.CheckResults).Sub(results.before); skew > time.Second || skew < -time.Second {
		t.Errorf("check_results cutoff off by %v", skew)
	}
	if skew := now.Add(-cfg.Incidents).Sub(incidents.before); skew > time.Second || skew < -time.Second {
		t.Errorf("incidents cutoff off by %v", skew)
	}
}

func TestCleaner_Cleanup_SkipsZeroWindows(t *testing.T) {
	results := &fakeResultRepo{}
	incidents := &fakeIncidentRepo{}
	cfg := config.RetentionConfig{
		CheckResults:    0, // disabled
		Incidents:       48 * time.Hour,
		CleanupInterval: time.Hour,
	}
	c := New(results, incidents, cfg, zerolog.Nop())

	c.cleanup(context.Background())

	if results.calls != 0 {
		t.Errorf("Expected check_results delete skipped, got %d calls", results.calls)
	}
	if incidents.calls != 1 {
		t.Errorf("Expected incidents delete called once, got %d", incidents.calls)
	}
}
