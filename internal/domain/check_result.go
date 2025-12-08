package domain

import "time"

// CheckStatus represents the status of a check
type CheckStatus string

const (
	CheckStatusSuccess CheckStatus = "success"
	CheckStatusFailure CheckStatus = "failure"
	CheckStatusWarning CheckStatus = "warning"
	CheckStatusUnknown CheckStatus = "unknown"
)

// CheckResult represents the result of a single check
type CheckResult struct {
	ID             int64                  `json:"id"`
	TargetID       string                 `json:"target_id"`
	Status         CheckStatus            `json:"status"`
	ResponseTimeMs int64                  `json:"response_time_ms"`
	StatusCode     int                    `json:"status_code,omitempty"`
	Error          string                 `json:"error,omitempty"`
	Message        string                 `json:"message,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CheckedAt      time.Time              `json:"checked_at"`
}

// IsHealthy returns true if the check was successful
func (cr *CheckResult) IsHealthy() bool {
	return cr.Status == CheckStatusSuccess
}

// CheckStats represents aggregated statistics for a target
type CheckStats struct {
	TargetID           string        `json:"target_id"`
	TotalChecks        int64         `json:"total_checks"`
	SuccessfulChecks   int64         `json:"successful_checks"`
	FailedChecks       int64         `json:"failed_checks"`
	UptimePercentage   float64       `json:"uptime_percentage"`
	AvgResponseTimeMs  float64       `json:"avg_response_time_ms"`
	MinResponseTimeMs  int64         `json:"min_response_time_ms"`
	MaxResponseTimeMs  int64         `json:"max_response_time_ms"`
	LastCheckAt        time.Time     `json:"last_check_at"`
	LastSuccessAt      time.Time     `json:"last_success_at,omitempty"`
	LastFailureAt      time.Time     `json:"last_failure_at,omitempty"`
	CurrentStatus      CheckStatus   `json:"current_status"`
	ConsecutiveFailures int          `json:"consecutive_failures"`
	Period             time.Duration `json:"period"`
}
