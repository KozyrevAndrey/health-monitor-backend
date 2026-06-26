package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping() error { return f.err }

// fakeScheduler embeds the interface so only IsRunning needs implementing.
type fakeScheduler struct {
	domain.Scheduler
	running bool
}

func (f fakeScheduler) IsRunning() bool { return f.running }

func doHealth(t *testing.T, s *Server) (int, HealthResponse) {
	t.Helper()
	rec := httptest.NewRecorder()
	s.handleHealth(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	var resp HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v (%s)", err, rec.Body.String())
	}
	return rec.Code, resp
}

func TestHandleHealth_AllOK(t *testing.T) {
	s := &Server{log: zerolog.Nop(), dbPinger: fakePinger{}, scheduler: fakeScheduler{running: true}}
	code, resp := doHealth(t, s)

	if code != http.StatusOK {
		t.Errorf("expected 200, got %d", code)
	}
	if resp.Status != "healthy" {
		t.Errorf("expected healthy, got %q", resp.Status)
	}
	if resp.Checks["database"].Status != "ok" || resp.Checks["scheduler"].Status != "ok" {
		t.Errorf("expected both checks ok, got %+v", resp.Checks)
	}
}

func TestHandleHealth_DBDown_Unhealthy(t *testing.T) {
	s := &Server{log: zerolog.Nop(), dbPinger: fakePinger{err: errors.New("connection refused")}, scheduler: fakeScheduler{running: true}}
	code, resp := doHealth(t, s)

	if code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", code)
	}
	if resp.Status != "unhealthy" {
		t.Errorf("expected unhealthy, got %q", resp.Status)
	}
	if resp.Checks["database"].Status != "error" {
		t.Errorf("expected database error, got %+v", resp.Checks["database"])
	}
}

func TestHandleHealth_SchedulerDown_Degraded(t *testing.T) {
	s := &Server{log: zerolog.Nop(), dbPinger: fakePinger{}, scheduler: fakeScheduler{running: false}}
	code, resp := doHealth(t, s)

	if code != http.StatusOK {
		t.Errorf("expected 200 for degraded, got %d", code)
	}
	if resp.Status != "degraded" {
		t.Errorf("expected degraded, got %q", resp.Status)
	}
}
