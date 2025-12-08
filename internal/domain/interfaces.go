package domain

import (
	"context"
	"time"
)

// Checker defines the interface for performing health checks
type Checker interface {
	// Check performs a health check on the target
	Check(ctx context.Context, target *Target) (*CheckResult, error)

	// Type returns the type of checker (http, tcp, dns, etc.)
	Type() TargetType

	// Validate validates the target configuration for this checker
	Validate(config map[string]interface{}) error
}

// TargetRepository defines the interface for target storage operations
type TargetRepository interface {
	// Create creates a new target
	Create(ctx context.Context, target *Target) error

	// Get retrieves a target by ID
	Get(ctx context.Context, id string) (*Target, error)

	// List retrieves all targets
	List(ctx context.Context) ([]*Target, error)

	// ListEnabled retrieves all enabled targets
	ListEnabled(ctx context.Context) ([]*Target, error)

	// Update updates an existing target
	Update(ctx context.Context, target *Target) error

	// Delete deletes a target by ID
	Delete(ctx context.Context, id string) error
}

// CheckResultRepository defines the interface for check result storage operations
type CheckResultRepository interface {
	// Save saves a check result
	Save(ctx context.Context, result *CheckResult) error

	// GetLatest retrieves the latest check result for a target
	GetLatest(ctx context.Context, targetID string) (*CheckResult, error)

	// GetHistory retrieves check history for a target
	GetHistory(ctx context.Context, targetID string, limit, offset int) ([]*CheckResult, error)

	// GetHistoryInRange retrieves check history within a time range
	GetHistoryInRange(ctx context.Context, targetID string, from, to time.Time) ([]*CheckResult, error)

	// GetStats retrieves aggregated statistics for a target
	GetStats(ctx context.Context, targetID string, period time.Duration) (*CheckStats, error)

	// DeleteOlderThan deletes check results older than the specified time
	DeleteOlderThan(ctx context.Context, before time.Time) (int64, error)
}

// IncidentRepository defines the interface for incident storage operations
type IncidentRepository interface {
	// Create creates a new incident
	Create(ctx context.Context, incident *Incident) error

	// Get retrieves an incident by ID
	Get(ctx context.Context, id int64) (*Incident, error)

	// GetOngoing retrieves ongoing incident for a target
	GetOngoing(ctx context.Context, targetID string) (*Incident, error)

	// List retrieves all incidents
	List(ctx context.Context, limit, offset int) ([]*Incident, error)

	// ListByTarget retrieves incidents for a specific target
	ListByTarget(ctx context.Context, targetID string, limit, offset int) ([]*Incident, error)

	// Update updates an incident
	Update(ctx context.Context, incident *Incident) error

	// Resolve marks an incident as resolved
	Resolve(ctx context.Context, id int64) error
}

// Notifier defines the interface for sending notifications
type Notifier interface {
	// Notify sends a notification
	Notify(ctx context.Context, alert *Alert) error

	// Type returns the type of notifier (webhook, email, telegram, etc.)
	Type() string

	// Validate validates the notifier configuration
	Validate(config map[string]interface{}) error
}

// Scheduler defines the interface for scheduling health checks
type Scheduler interface {
	// Start starts the scheduler
	Start(ctx context.Context) error

	// Stop stops the scheduler gracefully
	Stop() error

	// AddTarget adds a target to the scheduler
	AddTarget(target *Target) error

	// RemoveTarget removes a target from the scheduler
	RemoveTarget(targetID string) error

	// UpdateTarget updates a target in the scheduler
	UpdateTarget(target *Target) error

	// IsRunning returns true if the scheduler is running
	IsRunning() bool
}

// AlertManager defines the interface for managing alerts
type AlertManager interface {
	// ProcessCheckResult processes a check result and determines if alerts should be sent
	ProcessCheckResult(ctx context.Context, result *CheckResult) error

	// CreateAlert creates and sends an alert
	CreateAlert(ctx context.Context, alert *Alert) error

	// RegisterNotifier registers a notifier
	RegisterNotifier(notifier Notifier)

	// GetNotifier retrieves a notifier by type
	GetNotifier(notifierType string) (Notifier, error)
}

// CheckerRegistry defines the interface for managing checkers
type CheckerRegistry interface {
	// Register registers a new checker
	Register(checker Checker)

	// Get retrieves a checker by type
	Get(targetType TargetType) (Checker, error)

	// List returns all registered checker types
	List() []TargetType
}
