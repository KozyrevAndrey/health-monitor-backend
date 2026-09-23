package notifier

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"health-monitor/internal/domain"
)

// Every test here talks to an httptest server. Nothing may reach smspilot.ru:
// a real request places a phone call and costs money. The API key comes from a
// test-only environment variable, never from SMSPILOT_API_KEY, so a working key
// on a developer machine cannot leak into a test run.
const smspilotTestKeyEnv = "SMSPILOT_TEST_KEY"

func smspilotTestConfig(extra map[string]interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"phones":      []interface{}{"79000000000"},
		"api_key_env": smspilotTestKeyEnv,
	}
	for k, v := range extra {
		config[k] = v
	}
	return config
}

func newTestSMSPilot(t *testing.T, apiURL string, extra map[string]interface{}) *SMSPilotNotifier {
	t.Helper()

	t.Setenv(smspilotTestKeyEnv, "test-key")

	n, err := newSMSPilotNotifier(apiURL, smspilotTestConfig(extra), zerolog.Nop())
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}
	return n
}

func testAlert() *domain.Alert {
	return &domain.Alert{
		ID:         "alert-1",
		TargetID:   "target-1",
		TargetName: "Peptide server",
		Type:       domain.AlertTypeDown,
		Severity:   domain.AlertSeverityCritical,
		Message:    "Target Peptide server is DOWN",
		CreatedAt:  time.Now(),
	}
}

func TestSMSPilot_PlacesCallAndReturnsReceipt(t *testing.T) {
	var (
		mu    sync.Mutex
		forms []map[string]string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse form: %v", err)
		}

		form := make(map[string]string, len(r.PostForm))
		for k := range r.PostForm {
			form[k] = r.PostForm.Get(k)
		}

		mu.Lock()
		forms = append(forms, form)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"send":[{"server_id":"555","status":"1","phone":"79000000000","price":"3.5"}]}`))
	}))
	defer srv.Close()

	n := newTestSMSPilot(t, srv.URL, nil)

	receipt, err := n.NotifyWithReceipt(context.Background(), testAlert())
	if err != nil {
		t.Fatalf("NotifyWithReceipt failed: %v", err)
	}

	if receipt.ExternalID != "555" {
		t.Errorf("Expected server_id 555, got %q", receipt.ExternalID)
	}
	if receipt.ExternalStatus != "1" {
		t.Errorf("Expected status 1, got %q", receipt.ExternalStatus)
	}
	if receipt.Cost != "3.5" {
		t.Errorf("Expected cost 3.5, got %q", receipt.Cost)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(forms) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(forms))
	}

	form := forms[0]
	if form["to"] != "79000000000" {
		t.Errorf("Expected to=79000000000, got %q", form["to"])
	}
	if form["from"] != "GOLOS" {
		t.Errorf("Expected from=GOLOS, got %q", form["from"])
	}
	if form["apikey"] != "test-key" {
		t.Errorf("Expected the API key from the environment, got %q", form["apikey"])
	}
	if form["format"] != "json" {
		t.Errorf("Expected format=json, got %q", form["format"])
	}
	if !strings.Contains(form["send"], "Peptide server") {
		t.Errorf("Expected the target name in the message, got %q", form["send"])
	}
}

func TestSMSPilot_CustomTemplateAndVoice(t *testing.T) {
	var (
		mu   sync.Mutex
		text string
		from string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()

		mu.Lock()
		text = r.PostForm.Get("send")
		from = r.PostForm.Get("from")
		mu.Unlock()

		_, _ = w.Write([]byte(`{"send":[{"server_id":"1","status":"0"}]}`))
	}))
	defer srv.Close()

	n := newTestSMSPilot(t, srv.URL, map[string]interface{}{
		"message_template": "Сервис {{.TargetName}} не отвечает",
		"voice":            "golosm",
	})

	if _, err := n.NotifyWithReceipt(context.Background(), testAlert()); err != nil {
		t.Fatalf("NotifyWithReceipt failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if text != "Сервис Peptide server не отвечает" {
		t.Errorf("Unexpected rendered message: %q", text)
	}
	if from != "GOLOSM" {
		t.Errorf("Expected from=GOLOSM, got %q", from)
	}
}

func TestSMSPilot_CallsEveryPhone(t *testing.T) {
	var (
		mu    sync.Mutex
		calls []string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()

		mu.Lock()
		calls = append(calls, r.PostForm.Get("to"))
		mu.Unlock()

		_, _ = w.Write([]byte(`{"send":[{"server_id":"1","status":"1"}]}`))
	}))
	defer srv.Close()

	n := newTestSMSPilot(t, srv.URL, map[string]interface{}{
		"phones": []interface{}{"+7 900 000 00 01", "79000000002"},
	})

	receipt, err := n.NotifyWithReceipt(context.Background(), testAlert())
	if err != nil {
		t.Fatalf("NotifyWithReceipt failed: %v", err)
	}

	if receipt.ExternalID != "1,1" {
		t.Errorf("Expected a receipt per phone, got %q", receipt.ExternalID)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(calls) != 2 || calls[0] != "79000000001" || calls[1] != "79000000002" {
		t.Fatalf("Expected both numbers to be called and normalised, got %v", calls)
	}
}

// An explicit rejection means nothing was dialled: it must not be reported as
// an unknown result, otherwise the operator cannot tell the two apart.
func TestSMSPilot_ProviderErrorIsRejection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"error":{"code":"232","description":"whitelist","description_ru":"номер не в белом списке"}}`))
	}))
	defer srv.Close()

	n := newTestSMSPilot(t, srv.URL, nil)

	_, err := n.NotifyWithReceipt(context.Background(), testAlert())
	if err == nil {
		t.Fatal("Expected an error")
	}
	if errors.Is(err, domain.ErrUnknownResult) {
		t.Fatalf("Expected a rejection, got an unknown result: %v", err)
	}
	if !strings.Contains(err.Error(), "232") {
		t.Errorf("Expected the provider error code in the message, got %v", err)
	}
}

func TestSMSPilot_NegativeStatusIsRejection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"send":[{"server_id":"7","status":"-2","error_ru":"ошибка"}]}`))
	}))
	defer srv.Close()

	n := newTestSMSPilot(t, srv.URL, nil)

	_, err := n.NotifyWithReceipt(context.Background(), testAlert())
	if err == nil {
		t.Fatal("Expected an error")
	}
	if errors.Is(err, domain.ErrUnknownResult) {
		t.Fatalf("Expected a rejection, got an unknown result: %v", err)
	}
}

// Anything that leaves the outcome open must surface as ErrUnknownResult, which
// the dispatcher treats as terminal: the call may already have been placed.
func TestSMSPilot_UnknownResults(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name: "http error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		},
		{
			name: "unparseable body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("<html>gateway</html>"))
			},
		},
		{
			name: "no server_id",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"send":[{"status":"1"}]}`))
			},
		},
		{
			name: "ambiguous result",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"send":[]}`))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			n := newTestSMSPilot(t, srv.URL, nil)

			_, err := n.NotifyWithReceipt(context.Background(), testAlert())
			if !errors.Is(err, domain.ErrUnknownResult) {
				t.Fatalf("Expected ErrUnknownResult, got %v", err)
			}
		})
	}
}

func TestSMSPilot_ConfigValidation(t *testing.T) {
	t.Setenv(smspilotTestKeyEnv, "test-key")

	tests := []struct {
		name   string
		config map[string]interface{}
	}{
		{"no phones", map[string]interface{}{"api_key_env": smspilotTestKeyEnv}},
		{"landline", smspilotTestConfig(map[string]interface{}{"phones": []interface{}{"74951234567"}})},
		{"too short", smspilotTestConfig(map[string]interface{}{"phones": []interface{}{"7900000000"}})},
		{"bad voice", smspilotTestConfig(map[string]interface{}{"voice": "ALICE"})},
		{"broken template", smspilotTestConfig(map[string]interface{}{"message_template": "{{.TargetName"})},
		{"bad policy", smspilotTestConfig(map[string]interface{}{"repeat_interval": "soon"})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &SMSPilotNotifier{}
			if err := n.Validate(tt.config); err == nil {
				t.Fatal("Expected validation to fail")
			}

			// The API save path builds the notifier instead of calling
			// Validate, so the constructor must reject the same configs.
			if _, err := newSMSPilotNotifier("http://127.0.0.1:1", tt.config, zerolog.Nop()); err == nil {
				t.Fatal("Expected construction to fail")
			}
		})
	}
}

func TestSMSPilot_MissingAPIKeyIsRejected(t *testing.T) {
	t.Setenv(smspilotTestKeyEnv, "")

	n := &SMSPilotNotifier{}
	err := n.Validate(smspilotTestConfig(nil))
	if err == nil {
		t.Fatal("Expected validation to fail without an API key")
	}
	if !strings.Contains(err.Error(), smspilotTestKeyEnv) {
		t.Errorf("Expected the variable name in the error, got %v", err)
	}

	if _, err := newSMSPilotNotifier("http://127.0.0.1:1", smspilotTestConfig(nil), zerolog.Nop()); err == nil {
		t.Fatal("Expected construction to fail without an API key")
	}
}

func TestSMSPilot_ImplementsReceiptNotifier(t *testing.T) {
	var _ domain.ReceiptNotifier = (*SMSPilotNotifier)(nil)
	var _ domain.Notifier = (*SMSPilotNotifier)(nil)
}
