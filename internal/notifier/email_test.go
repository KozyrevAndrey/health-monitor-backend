package notifier

import (
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

func TestNewEmailNotifier(t *testing.T) {
	log := zerolog.Nop()

	tests := []struct {
		name      string
		config    map[string]interface{}
		shouldErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"smtp_host":     "smtp.gmail.com",
				"smtp_port":     587,
				"smtp_user":     "user@example.com",
				"smtp_password": "password",
				"from":          "alerts@example.com",
				"to":            []interface{}{"recipient@example.com"},
			},
			shouldErr: false,
		},
		{
			name: "valid config with defaults",
			config: map[string]interface{}{
				"smtp_host": "smtp.gmail.com",
				"from":      "alerts@example.com",
				"to":        []interface{}{"recipient@example.com"},
			},
			shouldErr: false,
		},
		{
			name: "missing smtp_host",
			config: map[string]interface{}{
				"from": "alerts@example.com",
				"to":   []interface{}{"recipient@example.com"},
			},
			shouldErr: true,
		},
		{
			name: "empty smtp_host",
			config: map[string]interface{}{
				"smtp_host": "",
				"from":      "alerts@example.com",
				"to":        []interface{}{"recipient@example.com"},
			},
			shouldErr: true,
		},
		{
			name: "missing from",
			config: map[string]interface{}{
				"smtp_host": "smtp.gmail.com",
				"to":        []interface{}{"recipient@example.com"},
			},
			shouldErr: true,
		},
		{
			name: "missing to",
			config: map[string]interface{}{
				"smtp_host": "smtp.gmail.com",
				"from":      "alerts@example.com",
			},
			shouldErr: true,
		},
		{
			name: "empty to array",
			config: map[string]interface{}{
				"smtp_host": "smtp.gmail.com",
				"from":      "alerts@example.com",
				"to":        []interface{}{},
			},
			shouldErr: true,
		},
		{
			name: "multiple recipients",
			config: map[string]interface{}{
				"smtp_host": "smtp.gmail.com",
				"from":      "alerts@example.com",
				"to":        []interface{}{"user1@example.com", "user2@example.com"},
			},
			shouldErr: false,
		},
		{
			name: "port as float64",
			config: map[string]interface{}{
				"smtp_host": "smtp.gmail.com",
				"smtp_port": float64(587),
				"from":      "alerts@example.com",
				"to":        []interface{}{"recipient@example.com"},
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier, err := NewEmailNotifier(tt.config, log)
			if tt.shouldErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if notifier == nil {
					t.Error("Expected notifier to be created")
				}
			}
		})
	}
}

func TestEmailNotifier_Type(t *testing.T) {
	log := zerolog.Nop()
	config := map[string]interface{}{
		"smtp_host": "smtp.gmail.com",
		"from":      "alerts@example.com",
		"to":        []interface{}{"recipient@example.com"},
	}

	notifier, err := NewEmailNotifier(config, log)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	if notifier.Type() != "email" {
		t.Errorf("Expected type 'email', got '%s'", notifier.Type())
	}
}

func TestEmailNotifier_Validate(t *testing.T) {
	log := zerolog.Nop()
	notifier := &EmailNotifier{log: log}

	tests := []struct {
		name      string
		config    map[string]interface{}
		shouldErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"smtp_host": "smtp.gmail.com",
				"from":      "alerts@example.com",
				"to":        []interface{}{"recipient@example.com"},
			},
			shouldErr: false,
		},
		{
			name: "invalid config",
			config: map[string]interface{}{
				"smtp_host": "",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := notifier.Validate(tt.config)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestEmailNotifier_BuildSubject(t *testing.T) {
	log := zerolog.Nop()
	config := map[string]interface{}{
		"smtp_host": "smtp.gmail.com",
		"from":      "alerts@example.com",
		"to":        []interface{}{"recipient@example.com"},
	}

	notifier, err := NewEmailNotifier(config, log)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	tests := []struct {
		name     string
		alert    *domain.Alert
		expected string
	}{
		{
			name: "critical alert",
			alert: &domain.Alert{
				TargetName: "Test Target",
				Type:       domain.AlertTypeDown,
				Severity:   domain.AlertSeverityCritical,
			},
			expected: "[CRITICAL] Test Target - DOWN",
		},
		{
			name: "warning alert",
			alert: &domain.Alert{
				TargetName: "Test Target",
				Type:       domain.AlertTypeSlowResponse,
				Severity:   domain.AlertSeverityWarning,
			},
			expected: "[WARNING] Test Target - SLOW_RESPONSE",
		},
		{
			name: "info alert",
			alert: &domain.Alert{
				TargetName: "Test Target",
				Type:       domain.AlertTypeUp,
				Severity:   domain.AlertSeverityInfo,
			},
			expected: "[INFO] Test Target - UP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subject := notifier.buildSubject(tt.alert)
			if subject != tt.expected {
				t.Errorf("Expected subject '%s', got '%s'", tt.expected, subject)
			}
		})
	}
}

func TestEmailNotifier_BuildPlainTextBody(t *testing.T) {
	log := zerolog.Nop()
	config := map[string]interface{}{
		"smtp_host": "smtp.gmail.com",
		"from":      "alerts@example.com",
		"to":        []interface{}{"recipient@example.com"},
	}

	notifier, err := NewEmailNotifier(config, log)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	alert := &domain.Alert{
		ID:          "test-alert-1",
		TargetID:    "test-target",
		TargetName:  "Test Target",
		Type:        domain.AlertTypeDown,
		Severity:    domain.AlertSeverityCritical,
		Message:     "Target is DOWN",
		Description: "Service unavailable",
		CreatedAt:   time.Now(),
		Metadata: map[string]interface{}{
			"status_code":      500,
			"response_time_ms": 1000,
		},
	}

	body := notifier.buildPlainTextBody(alert)

	if !strings.Contains(body, "Target is DOWN") {
		t.Error("Body should contain alert message")
	}

	if !strings.Contains(body, "Test Target") {
		t.Error("Body should contain target name")
	}

	if !strings.Contains(body, "down") {
		t.Error("Body should contain alert type")
	}

	if !strings.Contains(body, "critical") {
		t.Error("Body should contain severity")
	}

	if !strings.Contains(body, "Service unavailable") {
		t.Error("Body should contain description")
	}

	if !strings.Contains(body, "🔴") {
		t.Error("Body should contain DOWN icon")
	}

	if !strings.Contains(body, "Health Monitor") {
		t.Error("Body should contain footer")
	}
}

func TestEmailNotifier_BuildHTMLBody(t *testing.T) {
	log := zerolog.Nop()
	config := map[string]interface{}{
		"smtp_host": "smtp.gmail.com",
		"from":      "alerts@example.com",
		"to":        []interface{}{"recipient@example.com"},
	}

	notifier, err := NewEmailNotifier(config, log)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	alert := &domain.Alert{
		ID:          "test-alert-1",
		TargetID:    "test-target",
		TargetName:  "Test Target",
		Type:        domain.AlertTypeDown,
		Severity:    domain.AlertSeverityCritical,
		Message:     "Target is DOWN",
		Description: "Service unavailable",
		CreatedAt:   time.Now(),
		Metadata: map[string]interface{}{
			"status_code":      500,
			"response_time_ms": 1000,
		},
	}

	body := notifier.buildHTMLBody(alert)

	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Body should be valid HTML")
	}

	if !strings.Contains(body, "Target is DOWN") {
		t.Error("Body should contain alert message")
	}

	if !strings.Contains(body, "Test Target") {
		t.Error("Body should contain target name")
	}

	if !strings.Contains(body, "critical") {
		t.Error("Body should contain severity")
	}

	if !strings.Contains(body, "Service unavailable") {
		t.Error("Body should contain description")
	}

	if !strings.Contains(body, "#dc3545") {
		t.Error("Body should contain critical color (red)")
	}

	if !strings.Contains(body, "🔴") {
		t.Error("Body should contain DOWN icon")
	}

	if !strings.Contains(body, "Health Monitor") {
		t.Error("Body should contain footer")
	}
}

func TestEmailNotifier_GetSeverityColor(t *testing.T) {
	log := zerolog.Nop()
	notifier := &EmailNotifier{log: log}

	tests := []struct {
		severity domain.AlertSeverity
		expected string
	}{
		{domain.AlertSeverityCritical, "#dc3545"},
		{domain.AlertSeverityWarning, "#ffc107"},
		{domain.AlertSeverityInfo, "#17a2b8"},
	}

	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			color := notifier.getSeverityColor(tt.severity)
			if color != tt.expected {
				t.Errorf("Expected color '%s', got '%s'", tt.expected, color)
			}
		})
	}
}

func TestEmailNotifier_GetAlertIcon(t *testing.T) {
	log := zerolog.Nop()
	notifier := &EmailNotifier{log: log}

	tests := []struct {
		alertType domain.AlertType
		expected  string
	}{
		{domain.AlertTypeDown, "🔴"},
		{domain.AlertTypeUp, "🟢"},
		{domain.AlertTypeSlowResponse, "🐌"},
		{domain.AlertTypeSSLExpiring, "🔐"},
		{domain.AlertTypeConsecutiveFail, "⚠️"},
	}

	for _, tt := range tests {
		t.Run(string(tt.alertType), func(t *testing.T) {
			icon := notifier.getAlertIcon(tt.alertType)
			if icon != tt.expected {
				t.Errorf("Expected icon '%s', got '%s'", tt.expected, icon)
			}
		})
	}
}

func TestEmailNotifier_BuildEmailMessage(t *testing.T) {
	log := zerolog.Nop()
	config := map[string]interface{}{
		"smtp_host": "smtp.gmail.com",
		"from":      "alerts@example.com",
		"to":        []interface{}{"recipient@example.com"},
	}

	notifier, err := NewEmailNotifier(config, log)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	subject := "Test Subject"
	htmlBody := "<html><body>HTML Body</body></html>"
	plainBody := "Plain Text Body"

	msg := notifier.buildEmailMessage(subject, htmlBody, plainBody)

	if !strings.Contains(msg, "From: alerts@example.com") {
		t.Error("Message should contain From header")
	}

	if !strings.Contains(msg, "To: recipient@example.com") {
		t.Error("Message should contain To header")
	}

	if !strings.Contains(msg, "Subject: Test Subject") {
		t.Error("Message should contain Subject header")
	}

	if !strings.Contains(msg, "MIME-Version: 1.0") {
		t.Error("Message should contain MIME-Version header")
	}

	if !strings.Contains(msg, "multipart/alternative") {
		t.Error("Message should be multipart/alternative")
	}

	if !strings.Contains(msg, "text/plain") {
		t.Error("Message should contain plain text part")
	}

	if !strings.Contains(msg, "text/html") {
		t.Error("Message should contain HTML part")
	}

	if !strings.Contains(msg, "Plain Text Body") {
		t.Error("Message should contain plain text body")
	}

	if !strings.Contains(msg, "HTML Body") {
		t.Error("Message should contain HTML body")
	}
}
