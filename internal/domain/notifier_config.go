package domain

import "time"

// NotifierConfig represents a notification channel configuration stored in the database
type NotifierConfig struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Enabled   bool                   `json:"enabled"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
