package checker

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"health-monitor/internal/domain"
)

func TestTCPChecker_Type(t *testing.T) {
	if NewTCPChecker().Type() != domain.TargetTypeTCP {
		t.Errorf("Expected type %q", domain.TargetTypeTCP)
	}
}

func TestParseTCPConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    map[string]interface{}
		shouldErr bool
	}{
		{"valid float port", map[string]interface{}{"host": "example.com", "port": float64(443)}, false},
		{"valid int port", map[string]interface{}{"host": "example.com", "port": 80}, false},
		{"missing host", map[string]interface{}{"port": float64(80)}, true},
		{"empty host", map[string]interface{}{"host": "", "port": float64(80)}, true},
		{"missing port", map[string]interface{}{"host": "example.com"}, true},
		{"port too low", map[string]interface{}{"host": "example.com", "port": float64(0)}, true},
		{"port too high", map[string]interface{}{"host": "example.com", "port": float64(70000)}, true},
		{"port wrong type", map[string]interface{}{"host": "example.com", "port": "80"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseTCPConfig(tt.config)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTCPChecker_Check_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()

	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)

	checker := NewTCPChecker()
	target := &domain.Target{
		ID:      "tcp-ok",
		Type:    domain.TargetTypeTCP,
		Timeout: 5 * time.Second,
		Config:  map[string]interface{}{"host": "127.0.0.1", "port": float64(port)},
	}

	result, err := checker.Check(context.Background(), target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Status != domain.CheckStatusSuccess {
		t.Errorf("Expected success, got %s (%s)", result.Status, result.Message)
	}
	if result.Metadata["port"] != port {
		t.Errorf("Expected port metadata %d, got %v", port, result.Metadata["port"])
	}
}

func TestTCPChecker_Check_Failure(t *testing.T) {
	// Bind then close to obtain a port that is (almost certainly) not listening.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	ln.Close()

	checker := NewTCPChecker()
	target := &domain.Target{
		ID:      "tcp-fail",
		Type:    domain.TargetTypeTCP,
		Timeout: 2 * time.Second,
		Config:  map[string]interface{}{"host": "127.0.0.1", "port": float64(port)},
	}

	result, err := checker.Check(context.Background(), target)
	if err != nil {
		t.Fatalf("Check should not return error, got: %v", err)
	}
	if result.Status != domain.CheckStatusFailure {
		t.Errorf("Expected failure, got %s", result.Status)
	}
	if result.Error == "" {
		t.Error("Expected error message on failure")
	}
}
