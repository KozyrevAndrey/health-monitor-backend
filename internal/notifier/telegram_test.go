package notifier

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

func TestNewTelegramNotifier(t *testing.T) {
	log := zerolog.Nop()

	tests := []struct {
		name      string
		config    map[string]interface{}
		shouldErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
				"chat_ids":  []interface{}{"123456789"},
			},
			shouldErr: false,
		},
		{
			name: "missing bot_token",
			config: map[string]interface{}{
				"chat_ids": []interface{}{"123456789"},
			},
			shouldErr: true,
		},
		{
			name: "empty bot_token",
			config: map[string]interface{}{
				"bot_token": "",
				"chat_ids":  []interface{}{"123456789"},
			},
			shouldErr: true,
		},
		{
			name: "missing chat_ids",
			config: map[string]interface{}{
				"bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			},
			shouldErr: true,
		},
		{
			name: "empty chat_ids",
			config: map[string]interface{}{
				"bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
				"chat_ids":  []interface{}{},
			},
			shouldErr: true,
		},
		{
			name: "multiple chat_ids",
			config: map[string]interface{}{
				"bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
				"chat_ids":  []interface{}{"123456789", "987654321"},
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier, err := NewTelegramNotifier(tt.config, log)
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

func TestTelegramNotifier_Type(t *testing.T) {
	log := zerolog.Nop()
	config := map[string]interface{}{
		"bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		"chat_ids":  []interface{}{"123456789"},
	}

	notifier, err := NewTelegramNotifier(config, log)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	if notifier.Type() != "telegram" {
		t.Errorf("Expected type 'telegram', got '%s'", notifier.Type())
	}
}

func TestTelegramNotifier_Validate(t *testing.T) {
	log := zerolog.Nop()
	notifier := &TelegramNotifier{log: log}

	tests := []struct {
		name      string
		config    map[string]interface{}
		shouldErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
				"chat_ids":  []interface{}{"123456789"},
			},
			shouldErr: false,
		},
		{
			name: "invalid config",
			config: map[string]interface{}{
				"bot_token": "",
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

func TestTelegramNotifier_FormatMessage(t *testing.T) {
	log := zerolog.Nop()
	config := map[string]interface{}{
		"bot_token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		"chat_ids":  []interface{}{"123456789"},
	}

	notifier, err := NewTelegramNotifier(config, log)
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

	message := notifier.formatMessage(alert)

	if !strings.Contains(message, "Target is DOWN") {
		t.Error("Message should contain alert message")
	}

	if !strings.Contains(message, "Test Target") {
		t.Error("Message should contain target name")
	}

	if !strings.Contains(message, "down") {
		t.Error("Message should contain alert type")
	}

	if !strings.Contains(message, "critical") {
		t.Error("Message should contain severity")
	}

	if !strings.Contains(message, "Service unavailable") {
		t.Error("Message should contain description")
	}

	if !strings.Contains(message, "🔴") {
		t.Error("Message should contain DOWN icon")
	}
}

func TestTelegramNotifier_GetAlertIcon(t *testing.T) {
	log := zerolog.Nop()
	notifier := &TelegramNotifier{log: log}

	tests := []struct {
		alertType domain.AlertType
		severity  domain.AlertSeverity
		expected  string
	}{
		{domain.AlertTypeDown, domain.AlertSeverityCritical, "🔴"},
		{domain.AlertTypeUp, domain.AlertSeverityInfo, "🟢"},
		{domain.AlertTypeSlowResponse, domain.AlertSeverityWarning, "🐌"},
		{domain.AlertTypeSSLExpiring, domain.AlertSeverityWarning, "🔐"},
		{domain.AlertTypeConsecutiveFail, domain.AlertSeverityCritical, "⚠️"},
	}

	for _, tt := range tests {
		t.Run(string(tt.alertType), func(t *testing.T) {
			icon := notifier.getAlertIcon(tt.alertType, tt.severity)
			if icon != tt.expected {
				t.Errorf("Expected icon '%s', got '%s'", tt.expected, icon)
			}
		})
	}
}

func TestTelegramNotifier_EscapeMarkdown(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello World"},
		{"Test_with_underscore", "Test\\_with\\_underscore"},
		{"Test*with*asterisk", "Test\\*with\\*asterisk"},
		{"Test[bracket]", "Test\\[bracket\\]"},
		{"Test (parenthesis)", "Test \\(parenthesis\\)"},
		{"Price: $100.00", "Price: $100\\.00"},
		{"Error!", "Error\\!"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestTelegramNotifier_Notify_Integration(t *testing.T) {
	t.Skip("Integration test - run manually with: go test -v -run TestTelegramNotifier_Notify_Integration")

	log := zerolog.Nop()

	config := map[string]interface{}{
		"bot_token": "7612093885:AAHIIiimaI9IkUlD0WPLlXqzVqI_SBRhWkg",
		"chat_ids":  []interface{}{"YOUR_CHAT_ID_HERE"},
	}

	notifier, err := NewTelegramNotifier(config, log)
	if err != nil {
		t.Fatalf("Failed to create notifier: %v", err)
	}

	alert := &domain.Alert{
		ID:          "test-alert-1",
		TargetID:    "test-target",
		TargetName:  "Health Monitor Test",
		Type:        domain.AlertTypeDown,
		Severity:    domain.AlertSeverityCritical,
		Message:     "Test Alert from Health Monitor",
		Description: "This is a test notification to verify Telegram integration",
		CreatedAt:   time.Now(),
		Metadata: map[string]interface{}{
			"test": true,
			"env":  "development",
		},
	}

	ctx := context.Background()
	err = notifier.Notify(ctx, alert)
	if err != nil {
		t.Fatalf("Failed to send notification: %v", err)
	}

	t.Log("Notification sent successfully!")
}
