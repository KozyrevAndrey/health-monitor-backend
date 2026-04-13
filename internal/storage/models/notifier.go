package models

import (
	"encoding/json"
	"time"

	"health-monitor/internal/domain"
)

// NotifierConfig represents the database model for notifier configurations
type NotifierConfig struct {
	ID        string    `gorm:"primaryKey"`
	Type      string    `gorm:"not null;index"`
	Enabled   bool      `gorm:"not null;default:true;index"`
	Config    string    `gorm:"not null"` // JSON string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name
func (NotifierConfig) TableName() string {
	return "notifier_configs"
}

// ToDomain converts database model to domain model
func (n *NotifierConfig) ToDomain() (*domain.NotifierConfig, error) {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(n.Config), &config); err != nil {
		return nil, err
	}

	return &domain.NotifierConfig{
		ID:        n.ID,
		Type:      n.Type,
		Enabled:   n.Enabled,
		Config:    config,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}, nil
}

// FromDomain converts domain model to database model
func (n *NotifierConfig) FromDomain(d *domain.NotifierConfig) error {
	configJSON, err := json.Marshal(d.Config)
	if err != nil {
		return err
	}

	n.ID = d.ID
	n.Type = d.Type
	n.Enabled = d.Enabled
	n.Config = string(configJSON)
	n.CreatedAt = d.CreatedAt
	n.UpdatedAt = d.UpdatedAt

	return nil
}
