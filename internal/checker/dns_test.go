package checker

import (
	"context"
	"testing"
	"time"

	"health-monitor/internal/domain"
)

func TestDNSChecker_Type(t *testing.T) {
	if NewDNSChecker().Type() != domain.TargetTypeDNS {
		t.Errorf("Expected type %q", domain.TargetTypeDNS)
	}
}

func TestParseDNSConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    map[string]interface{}
		shouldErr bool
		wantType  string
	}{
		{"defaults to A", map[string]interface{}{"domain": "example.com"}, false, "A"},
		{"lowercase record type normalized", map[string]interface{}{"domain": "example.com", "record_type": "mx"}, false, "MX"},
		{"missing domain", map[string]interface{}{"record_type": "A"}, true, ""},
		{"empty domain", map[string]interface{}{"domain": ""}, true, ""},
		{"unsupported record type", map[string]interface{}{"domain": "example.com", "record_type": "SRV"}, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseDNSConfig(tt.config)
			if tt.shouldErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if cfg.RecordType != tt.wantType {
				t.Errorf("Expected record type %q, got %q", tt.wantType, cfg.RecordType)
			}
		})
	}
}

func TestParseDNSConfig_ExpectedSlices(t *testing.T) {
	cfg, err := parseDNSConfig(map[string]interface{}{
		"domain":        "example.com",
		"expected_ips":  []interface{}{"1.2.3.4", "5.6.7.8"},
		"expect_values": []interface{}{"ns1.example.com"},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(cfg.ExpectedIPs) != 2 || len(cfg.ExpectValues) != 1 {
		t.Errorf("Expected 2 IPs and 1 value, got %v / %v", cfg.ExpectedIPs, cfg.ExpectValues)
	}
}

func TestMissingValues(t *testing.T) {
	got := []string{"1.2.3.4", "mx1.example.com (pref 10)", "NS1.Example.com"}
	tests := []struct {
		name     string
		expected []string
		want     int
	}{
		{"exact ip present", []string{"1.2.3.4"}, 0},
		{"mx host matches prefix", []string{"mx1.example.com"}, 0},
		{"case and trailing dot tolerant", []string{"ns1.example.com."}, 0},
		{"missing one", []string{"9.9.9.9"}, 1},
		{"empty expected", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if n := len(missingValues(tt.expected, got)); n != tt.want {
				t.Errorf("Expected %d missing, got %d", tt.want, n)
			}
		})
	}
}

func TestDNSChecker_Check_LocalhostA(t *testing.T) {
	checker := NewDNSChecker()
	target := &domain.Target{
		ID:      "dns-localhost",
		Type:    domain.TargetTypeDNS,
		Timeout: 5 * time.Second,
		Config: map[string]interface{}{
			"domain":       "localhost",
			"record_type":  "A",
			"expected_ips": []interface{}{"127.0.0.1"},
		},
	}

	result, err := checker.Check(context.Background(), target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Status != domain.CheckStatusSuccess {
		t.Errorf("Expected success resolving localhost, got %s (%s)", result.Status, result.Message)
	}
}

func TestDNSChecker_Check_NXDomain(t *testing.T) {
	checker := NewDNSChecker()
	target := &domain.Target{
		ID:      "dns-nx",
		Type:    domain.TargetTypeDNS,
		Timeout: 5 * time.Second,
		Config: map[string]interface{}{
			// .invalid is a reserved TLD guaranteed never to resolve.
			"domain":      "this-does-not-exist.invalid",
			"record_type": "A",
		},
	}

	result, err := checker.Check(context.Background(), target)
	if err != nil {
		t.Fatalf("Check should not return error, got: %v", err)
	}
	if result.Status != domain.CheckStatusFailure {
		t.Errorf("Expected failure for .invalid domain, got %s", result.Status)
	}
}
