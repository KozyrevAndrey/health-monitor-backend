package models

import (
	"encoding/json"
	"time"

	"health-monitor/internal/domain"
)

// CheckResult represents the database model for check results
type CheckResult struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	TargetID       string    `gorm:"not null;index:idx_target_checked,priority:1"`
	Status         string    `gorm:"not null;index"`
	ResponseTimeMs int64     `gorm:"not null"`
	StatusCode     *int      `gorm:"default:null"`
	Error          string    `gorm:"type:text"`
	Message        string    `gorm:"type:text"`
	Metadata       string    `gorm:"type:text"` // JSON
	CheckedAt      time.Time `gorm:"not null;index;index:idx_target_checked,priority:2,sort:desc"`
}

// TableName specifies the table name
func (CheckResult) TableName() string {
	return "check_results"
}

// ToDomain converts database model to domain model
func (cr *CheckResult) ToDomain() (*domain.CheckResult, error) {
	// Parse metadata
	var metadata map[string]interface{}
	if cr.Metadata != "" {
		if err := json.Unmarshal([]byte(cr.Metadata), &metadata); err != nil {
			return nil, err
		}
	}

	var statusCode int
	if cr.StatusCode != nil {
		statusCode = *cr.StatusCode
	}

	return &domain.CheckResult{
		ID:             cr.ID,
		TargetID:       cr.TargetID,
		Status:         domain.CheckStatus(cr.Status),
		ResponseTimeMs: cr.ResponseTimeMs,
		StatusCode:     statusCode,
		Error:          cr.Error,
		Message:        cr.Message,
		Metadata:       metadata,
		CheckedAt:      cr.CheckedAt,
	}, nil
}

// FromDomain converts domain model to database model
func (cr *CheckResult) FromDomain(dcr *domain.CheckResult) error {
	// Marshal metadata
	var metadataJSON []byte
	var err error
	if len(dcr.Metadata) > 0 {
		metadataJSON, err = json.Marshal(dcr.Metadata)
		if err != nil {
			return err
		}
	}

	cr.ID = dcr.ID
	cr.TargetID = dcr.TargetID
	cr.Status = string(dcr.Status)
	cr.ResponseTimeMs = dcr.ResponseTimeMs
	if dcr.StatusCode > 0 {
		cr.StatusCode = &dcr.StatusCode
	}
	cr.Error = dcr.Error
	cr.Message = dcr.Message
	cr.Metadata = string(metadataJSON)
	cr.CheckedAt = dcr.CheckedAt

	return nil
}
