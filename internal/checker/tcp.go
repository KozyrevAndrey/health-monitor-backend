package checker

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"health-monitor/internal/domain"
)

// TCPChecker verifies that a TCP port is reachable and measures the time taken
// to establish the connection.
type TCPChecker struct{}

func NewTCPChecker() *TCPChecker {
	return &TCPChecker{}
}

func (c *TCPChecker) Check(ctx context.Context, target *domain.Target) (*domain.CheckResult, error) {
	tcpConfig, err := parseTCPConfig(target.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid TCP config: %w", err)
	}

	result := &domain.CheckResult{
		TargetID:  target.ID,
		CheckedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	address := net.JoinHostPort(tcpConfig.Host, strconv.Itoa(tcpConfig.Port))
	result.Metadata["host"] = tcpConfig.Host
	result.Metadata["port"] = tcpConfig.Port
	result.Metadata["address"] = address

	checkCtx, cancel := context.WithTimeout(ctx, target.Timeout)
	defer cancel()

	var dialer net.Dialer
	startTime := time.Now()
	conn, err := dialer.DialContext(checkCtx, "tcp", address)
	responseTime := time.Since(startTime)
	result.ResponseTimeMs = responseTime.Milliseconds()

	if err != nil {
		result.Status = domain.CheckStatusFailure
		result.Error = err.Error()
		result.Message = fmt.Sprintf("TCP connection to %s failed: %s", address, err.Error())
		return result, nil
	}
	_ = conn.Close()

	result.Status = domain.CheckStatusSuccess
	result.Message = fmt.Sprintf("TCP connection to %s established in %dms", address, responseTime.Milliseconds())

	return result, nil
}

func (c *TCPChecker) Type() domain.TargetType {
	return domain.TargetTypeTCP
}

func (c *TCPChecker) Validate(config map[string]interface{}) error {
	_, err := parseTCPConfig(config)
	return err
}

func parseTCPConfig(config map[string]interface{}) (*domain.TCPConfig, error) {
	tcpConfig := &domain.TCPConfig{}

	host, ok := config["host"].(string)
	if !ok || host == "" {
		return nil, fmt.Errorf("host is required")
	}
	tcpConfig.Host = host

	port, err := configInt(config["port"])
	if err != nil {
		return nil, fmt.Errorf("port is required and must be a number")
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535")
	}
	tcpConfig.Port = port

	return tcpConfig, nil
}

// configInt coerces a JSON-decoded config value (float64) or an int into an int.
func configInt(raw interface{}) (int, error) {
	switch v := raw.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("not a number")
	}
}
