package models

import (
	"time"

	"health-monitor/internal/domain"
)

// NotificationSend represents the database model for notification sends.
// The composite unique index is what guarantees one notification per
// (incident, notifier, sequence), across restarts and concurrent checks.
type NotificationSend struct {
	ID             int64      `gorm:"primaryKey;autoIncrement"`
	IncidentID     int64      `gorm:"not null;uniqueIndex:idx_send_dedupe,priority:1"`
	NotifierID     string     `gorm:"not null;uniqueIndex:idx_send_dedupe,priority:2"`
	Sequence       int        `gorm:"not null;uniqueIndex:idx_send_dedupe,priority:3"`
	Status         string     `gorm:"not null;index"`
	ExternalID     string     `gorm:"default:''"`
	ExternalStatus string     `gorm:"default:''"`
	Cost           string     `gorm:"default:''"`
	Error          string     `gorm:"type:text"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;index"`
	CompletedAt    *time.Time `gorm:"default:null"`
}

// TableName specifies the table name
func (NotificationSend) TableName() string {
	return "notification_sends"
}

// ToDomain converts database model to domain model
func (n *NotificationSend) ToDomain() *domain.NotificationSend {
	return &domain.NotificationSend{
		ID:             n.ID,
		IncidentID:     n.IncidentID,
		NotifierID:     n.NotifierID,
		Sequence:       n.Sequence,
		Status:         domain.SendStatus(n.Status),
		ExternalID:     n.ExternalID,
		ExternalStatus: n.ExternalStatus,
		Cost:           n.Cost,
		Error:          n.Error,
		CreatedAt:      n.CreatedAt,
		CompletedAt:    n.CompletedAt,
	}
}

// FromDomain converts domain model to database model
func (n *NotificationSend) FromDomain(d *domain.NotificationSend) {
	n.ID = d.ID
	n.IncidentID = d.IncidentID
	n.NotifierID = d.NotifierID
	n.Sequence = d.Sequence
	n.Status = string(d.Status)
	n.ExternalID = d.ExternalID
	n.ExternalStatus = d.ExternalStatus
	n.Cost = d.Cost
	n.Error = d.Error
	n.CompletedAt = d.CompletedAt
}
