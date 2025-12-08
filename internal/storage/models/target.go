package models

import (
	"encoding/json"
	"time"

	"health-monitor/internal/domain"
)

// Target represents the database model for targets
type Target struct {
	ID          string    `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Type        string    `gorm:"not null;index"`
	Config      string    `gorm:"not null"` // JSON string
	Interval    int64     `gorm:"not null"` // nanoseconds
	Timeout     int64     `gorm:"not null"` // nanoseconds
	Enabled     bool      `gorm:"not null;default:true;index"`
	Tags        string    `gorm:"type:text"` // JSON array
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name
func (Target) TableName() string {
	return "targets"
}

// ToDomain converts database model to domain model
func (t *Target) ToDomain() (*domain.Target, error) {
	// Parse config
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(t.Config), &config); err != nil {
		return nil, err
	}

	// Parse tags
	var tags []string
	if t.Tags != "" {
		if err := json.Unmarshal([]byte(t.Tags), &tags); err != nil {
			return nil, err
		}
	}

	return &domain.Target{
		ID:          t.ID,
		Name:        t.Name,
		Type:        domain.TargetType(t.Type),
		Config:      config,
		Interval:    time.Duration(t.Interval),
		Timeout:     time.Duration(t.Timeout),
		Enabled:     t.Enabled,
		Tags:        tags,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}, nil
}

// FromDomain converts domain model to database model
func (t *Target) FromDomain(dt *domain.Target) error {
	// Marshal config
	configJSON, err := json.Marshal(dt.Config)
	if err != nil {
		return err
	}

	// Marshal tags
	var tagsJSON []byte
	if len(dt.Tags) > 0 {
		tagsJSON, err = json.Marshal(dt.Tags)
		if err != nil {
			return err
		}
	}

	t.ID = dt.ID
	t.Name = dt.Name
	t.Type = string(dt.Type)
	t.Config = string(configJSON)
	t.Interval = int64(dt.Interval)
	t.Timeout = int64(dt.Timeout)
	t.Enabled = dt.Enabled
	t.Tags = string(tagsJSON)
	t.Description = dt.Description
	t.CreatedAt = dt.CreatedAt
	t.UpdatedAt = dt.UpdatedAt

	return nil
}
