package notifier

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

func sampleAlert() *domain.Alert {
	return &domain.Alert{
		ID:         "alert-1",
		TargetID:   "target-1",
		TargetName: "Test Target",
		Type:       domain.AlertTypeDown,
		Severity:   domain.AlertSeverityCritical,
		Message:    "Target is DOWN",
		CreatedAt:  time.Now(),
	}
}

func TestParseWebhookConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    map[string]interface{}
		shouldErr bool
	}{
		{
			name:   "minimal valid",
			config: map[string]interface{}{"url": "https://example.com/hook"},
		},
		{
			name:      "missing url",
			config:    map[string]interface{}{},
			shouldErr: true,
		},
		{
			name:      "non-http scheme",
			config:    map[string]interface{}{"url": "ftp://example.com"},
			shouldErr: true,
		},
		{
			name: "valid full",
			config: map[string]interface{}{
				"url":         "https://example.com/hook",
				"method":      "put",
				"headers":     map[string]interface{}{"Authorization": "Bearer x"},
				"payload":     `{"text":"{{.Message}}"}`,
				"timeout":     "5s",
				"max_retries": float64(2),
			},
		},
		{
			name:      "bad timeout",
			config:    map[string]interface{}{"url": "https://example.com", "timeout": "nope"},
			shouldErr: true,
		},
		{
			name:      "negative retries",
			config:    map[string]interface{}{"url": "https://example.com", "max_retries": float64(-1)},
			shouldErr: true,
		},
		{
			name:      "bad payload template",
			config:    map[string]interface{}{"url": "https://example.com", "payload": "{{.Message"},
			shouldErr: true,
		},
		{
			name:      "non-string header value",
			config:    map[string]interface{}{"url": "https://example.com", "headers": map[string]interface{}{"X": 1}},
			shouldErr: true,
		},
		{
			name:      "unsupported proxy",
			config:    map[string]interface{}{"url": "https://example.com", "proxy_url": "ftp://h:1"},
			shouldErr: true,
		},
	}

	log := zerolog.Nop()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWebhookNotifier(tt.config, log)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestParseWebhookConfig_Defaults(t *testing.T) {
	cfg, err := parseWebhookConfig(map[string]interface{}{"url": "https://example.com/hook"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cfg.Method != http.MethodPost {
		t.Errorf("Expected default method POST, got %s", cfg.Method)
	}
	if cfg.Timeout != webhookDefaultTimeout {
		t.Errorf("Expected default timeout %v, got %v", webhookDefaultTimeout, cfg.Timeout)
	}
	if cfg.MaxRetries != webhookDefaultMaxRetries {
		t.Errorf("Expected default retries %d, got %d", webhookDefaultMaxRetries, cfg.MaxRetries)
	}
}

func TestWebhookNotifier_Type(t *testing.T) {
	n, err := NewWebhookNotifier(map[string]interface{}{"url": "https://example.com"}, zerolog.Nop())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if n.Type() != "webhook" {
		t.Errorf("Expected type 'webhook', got '%s'", n.Type())
	}
}

func TestWebhookNotifier_DefaultPayloadAndHeaders(t *testing.T) {
	var gotBody string
	var gotCT string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		gotCT = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n, err := NewWebhookNotifier(map[string]interface{}{"url": srv.URL}, zerolog.Nop())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err := n.Notify(context.Background(), sampleAlert()); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}
	if gotCT != "application/json" {
		t.Errorf("Expected default Content-Type application/json, got %q", gotCT)
	}
	if !strings.Contains(gotBody, "\"target_name\":\"Test Target\"") {
		t.Errorf("Expected default JSON payload with alert fields, got %s", gotBody)
	}
}

func TestWebhookNotifier_TemplatePayload(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n, err := NewWebhookNotifier(map[string]interface{}{
		"url":     srv.URL,
		"payload": `{"text":"{{.Message}} on {{.TargetName}}"}`,
	}, zerolog.Nop())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err := n.Notify(context.Background(), sampleAlert()); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}
	want := `{"text":"Target is DOWN on Test Target"}`
	if gotBody != want {
		t.Errorf("Expected body %q, got %q", want, gotBody)
	}
}

func TestWebhookNotifier_RetriesOn5xxThenSucceeds(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n, err := NewWebhookNotifier(map[string]interface{}{
		"url":         srv.URL,
		"max_retries": float64(3),
	}, zerolog.Nop())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Shorten backoff impact: base is 500ms, 2 retries => ~1.5s total. Acceptable for test.
	if err := n.Notify(context.Background(), sampleAlert()); err != nil {
		t.Fatalf("Notify should succeed after retries: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Errorf("Expected 3 calls, got %d", got)
	}
}

func TestWebhookNotifier_NoRetryOn4xx(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	n, err := NewWebhookNotifier(map[string]interface{}{
		"url":         srv.URL,
		"max_retries": float64(3),
	}, zerolog.Nop())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err := n.Notify(context.Background(), sampleAlert()); err == nil {
		t.Fatal("Expected error on 4xx")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("Expected exactly 1 call (no retry on 4xx), got %d", got)
	}
}
