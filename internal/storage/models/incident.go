package models

import (
	"time"

	"health-monitor/internal/domain"
)

// Incident represents the database model for incidents
type Incident struct {
	ID                   int64      `gorm:"primaryKey;autoIncrement"`
	TargetID             string     `gorm:"not null;index:idx_target_status,priority:1"`
	TargetName           string     `gorm:"not null"`
	Status               string     `gorm:"not null;index;index:idx_target_status,priority:2"`
	StartedAt            time.Time  `gorm:"not null;index:idx_started_at,sort:desc"`
	ResolvedAt           *time.Time `gorm:"default:null"`
	Duration             *int64     `gorm:"default:null"` // nanoseconds
	FailureCount         int        `gorm:"not null;default:0"`
	LastError            string     `gorm:"type:text"`
	AlertsSent           int        `gorm:"not null;default:0"`
	Severity             string     `gorm:"not null"`
	FirstCheckResultID   *int64     `gorm:"default:null"`
	LastCheckResultID    *int64     `gorm:"default:null"`
}

// TableName specifies the table name
func (Incident) TableName() string {
	return "incidents"
}

// ToDomain converts database model to domain model
func (i *Incident) ToDomain() *domain.Incident {
	var duration time.Duration
	if i.Duration != nil {
		duration = time.Duration(*i.Duration)
	}

	incident := &domain.Incident{
		ID:           i.ID,
		TargetID:     i.TargetID,
		TargetName:   i.TargetName,
		Status:       domain.IncidentStatus(i.Status),
		StartedAt:    i.StartedAt,
		ResolvedAt:   i.ResolvedAt,
		Duration:     duration,
		FailureCount: i.FailureCount,
		LastError:    i.LastError,
		AlertsSent:   i.AlertsSent,
		Severity:     domain.AlertSeverity(i.Severity),
	}

	// Note: FirstCheckResult and LastCheckResult are not populated here
	// They can be loaded separately if needed

	return incident
}

// FromDomain converts domain model to database model
func (i *Incident) FromDomain(di *domain.Incident) {
	i.ID = di.ID
	i.TargetID = di.TargetID
	i.TargetName = di.TargetName
	i.Status = string(di.Status)
	i.StartedAt = di.StartedAt
	i.ResolvedAt = di.ResolvedAt
	if di.Duration > 0 {
		duration := int64(di.Duration)
		i.Duration = &duration
	}
	i.FailureCount = di.FailureCount
	i.LastError = di.LastError
	i.AlertsSent = di.AlertsSent
	i.Severity = string(di.Severity)

	// Store check result IDs if available
	if di.FirstCheckResult != nil && di.FirstCheckResult.ID > 0 {
		i.FirstCheckResultID = &di.FirstCheckResult.ID
	}
	if di.LastCheckResult != nil && di.LastCheckResult.ID > 0 {
		i.LastCheckResultID = &di.LastCheckResult.ID
	}
}
