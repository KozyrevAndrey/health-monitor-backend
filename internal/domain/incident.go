package domain

import "time"

// IncidentStatus represents the status of an incident
type IncidentStatus string

const (
	IncidentStatusOngoing  IncidentStatus = "ongoing"
	IncidentStatusResolved IncidentStatus = "resolved"
)

// Incident represents a period of downtime or issues
type Incident struct {
	ID               int64          `json:"id"`
	TargetID         string         `json:"target_id"`
	TargetName       string         `json:"target_name"`
	Status           IncidentStatus `json:"status"`
	StartedAt        time.Time      `json:"started_at"`
	ResolvedAt       *time.Time     `json:"resolved_at,omitempty"`
	Duration         time.Duration  `json:"duration"`
	FailureCount     int            `json:"failure_count"`
	LastError        string         `json:"last_error,omitempty"`
	AlertsSent       int            `json:"alerts_sent"`
	Severity         AlertSeverity  `json:"severity"`
	FirstCheckResult *CheckResult   `json:"first_check_result,omitempty"`
	LastCheckResult  *CheckResult   `json:"last_check_result,omitempty"`
}

// IsOngoing returns true if the incident is still ongoing
func (i *Incident) IsOngoing() bool {
	return i.Status == IncidentStatusOngoing
}

// Resolve marks the incident as resolved
func (i *Incident) Resolve() {
	now := time.Now()
	i.Status = IncidentStatusResolved
	i.ResolvedAt = &now
	i.Duration = now.Sub(i.StartedAt)
}
