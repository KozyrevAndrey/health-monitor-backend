package checker

import (
	"context"
	"testing"
	"time"

	"health-monitor/internal/domain"
)

func TestHTTPChecker_Check_Success(t *testing.T) {
	checker := NewHTTPChecker()

	target := &domain.Target{
		ID:      "test-alfabank",
		Name:    "Alfabank Test",
		Type:    domain.TargetTypeHTTP,
		Timeout: 10 * time.Second,
		Config: map[string]interface{}{
			"url":                  "https://alfabank.far-harbor.ru/",
			"method":               "GET",
			"expected_status_code": 200,
			"validate_ssl":         true,
			"check_ssl_expiry":     true,
		},
	}

	ctx := context.Background()
	result, err := checker.Check(ctx, target)

	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	t.Logf("Status: %s", result.Status)
	t.Logf("Message: %s", result.Message)
	t.Logf("Response Time: %dms", result.ResponseTimeMs)
	t.Logf("Status Code: %d", result.StatusCode)

	if result.Metadata != nil {
		if sslDays, ok := result.Metadata["ssl_expiry_days"]; ok {
			t.Logf("SSL Expiry Days: %v", sslDays)
		}
		if sslExpires, ok := result.Metadata["ssl_expires_at"]; ok {
			t.Logf("SSL Expires At: %v", sslExpires)
		}
	}

	if result.Status != domain.CheckStatusSuccess && result.Status != domain.CheckStatusWarning {
		t.Errorf("Expected success or warning, got: %s (error: %s)", result.Status, result.Error)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %d", result.StatusCode)
	}
}

func TestHTTPChecker_Check_InvalidURL(t *testing.T) {
	checker := NewHTTPChecker()

	target := &domain.Target{
		ID:      "test-invalid",
		Name:    "Invalid Test",
		Type:    domain.TargetTypeHTTP,
		Timeout: 5 * time.Second,
		Config: map[string]interface{}{
			"url":                  "https://this-domain-definitely-does-not-exist-12345.com",
			"method":               "GET",
			"expected_status_code": 200,
		},
	}

	ctx := context.Background()
	result, err := checker.Check(ctx, target)

	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}

	if result.Status != domain.CheckStatusFailure {
		t.Errorf("Expected failure status, got: %s", result.Status)
	}

	t.Logf("Error message: %s", result.Error)
}

func TestHTTPChecker_Check_WrongStatusCode(t *testing.T) {
	checker := NewHTTPChecker()

	target := &domain.Target{
		ID:      "test-wrong-code",
		Name:    "Wrong Status Code Test",
		Type:    domain.TargetTypeHTTP,
		Timeout: 10 * time.Second,
		Config: map[string]interface{}{
			"url":                  "https://alfabank.far-harbor.ru/",
			"method":               "GET",
			"expected_status_code": 404,
		},
	}

	ctx := context.Background()
	result, err := checker.Check(ctx, target)

	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if result.Status != domain.CheckStatusFailure {
		t.Errorf("Expected failure for wrong status code, got: %s", result.Status)
	}

	t.Logf("Message: %s", result.Message)
}

func TestHTTPChecker_Validate(t *testing.T) {
	checker := NewHTTPChecker()

	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"url":                  "https://example.com",
				"method":               "GET",
				"expected_status_code": 200,
			},
			wantErr: false,
		},
		{
			name:    "missing url",
			config:  map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "empty url",
			config: map[string]interface{}{
				"url": "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checker.Validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
