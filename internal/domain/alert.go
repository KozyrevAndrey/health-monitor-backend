package domain

import "time"

// AlertSeverity represents the severity of an alert
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
)

// AlertType represents the type of alert
type AlertType string

const (
	AlertTypeDown            AlertType = "down"
	AlertTypeUp              AlertType = "up"
	AlertTypeSlowResponse    AlertType = "slow_response"
	AlertTypeSSLExpiring     AlertType = "ssl_expiring"
	AlertTypeConsecutiveFail AlertType = "consecutive_fail"
)

// Alert represents an alert that needs to be sent
type Alert struct {
	ID          string                 `json:"id"`
	TargetID    string                 `json:"target_id"`
	TargetName  string                 `json:"target_name"`
	Type        AlertType              `json:"type"`
	Severity    AlertSeverity          `json:"severity"`
	Message     string                 `json:"message"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Resolved    bool                   `json:"resolved"`
}

// AlertRule defines when an alert should be triggered
type AlertRule struct {
	ID                  string        `json:"id"`
	Name                string        `json:"name"`
	Enabled             bool          `json:"enabled"`
	ConsecutiveFailures int           `json:"consecutive_failures,omitempty"`
	ResponseTimeMs      int           `json:"response_time_ms,omitempty"`
	SSLExpiryDays       int           `json:"ssl_expiry_days,omitempty"`
	Severity            AlertSeverity `json:"severity"`
	NotifierIDs         []string      `json:"notifier_ids"`
}
