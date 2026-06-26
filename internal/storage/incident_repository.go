package storage

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"health-monitor/internal/domain"
	"health-monitor/internal/storage/models"
)

// IncidentRepository implements domain.IncidentRepository
type IncidentRepository struct {
	db *gorm.DB
}

// NewIncidentRepository creates a new incident repository
func NewIncidentRepository(db *gorm.DB) *IncidentRepository {
	return &IncidentRepository{db: db}
}

// Create creates a new incident
func (r *IncidentRepository) Create(ctx context.Context, incident *domain.Incident) error {
	model := &models.Incident{}
	model.FromDomain(incident)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create incident: %w", err)
	}

	// Update domain model with generated ID
	incident.ID = model.ID

	return nil
}

// Get retrieves an incident by ID
func (r *IncidentRepository) Get(ctx context.Context, id int64) (*domain.Incident, error) {
	var model models.Incident
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("incident not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}

	return model.ToDomain(), nil
}

// GetOngoing retrieves ongoing incident for a target
func (r *IncidentRepository) GetOngoing(ctx context.Context, targetID string) (*domain.Incident, error) {
	var model models.Incident
	if err := r.db.WithContext(ctx).
		Where("target_id = ? AND status = ?", targetID, domain.IncidentStatusOngoing).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No ongoing incident is not an error
		}
		return nil, fmt.Errorf("failed to get ongoing incident: %w", err)
	}

	return model.ToDomain(), nil
}

// List retrieves all incidents
func (r *IncidentRepository) List(ctx context.Context, limit, offset int) ([]*domain.Incident, error) {
	var modelIncidents []models.Incident
	query := r.db.WithContext(ctx).Order("started_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&modelIncidents).Error; err != nil {
		return nil, fmt.Errorf("failed to list incidents: %w", err)
	}

	incidents := make([]*domain.Incident, 0, len(modelIncidents))
	for _, model := range modelIncidents {
		incidents = append(incidents, model.ToDomain())
	}

	return incidents, nil
}

// ListByTarget retrieves incidents for a specific target
func (r *IncidentRepository) ListByTarget(ctx context.Context, targetID string, limit, offset int) ([]*domain.Incident, error) {
	var modelIncidents []models.Incident
	query := r.db.WithContext(ctx).
		Where("target_id = ?", targetID).
		Order("started_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&modelIncidents).Error; err != nil {
		return nil, fmt.Errorf("failed to list incidents by target: %w", err)
	}

	incidents := make([]*domain.Incident, 0, len(modelIncidents))
	for _, model := range modelIncidents {
		incidents = append(incidents, model.ToDomain())
	}

	return incidents, nil
}

// Update updates an incident
func (r *IncidentRepository) Update(ctx context.Context, incident *domain.Incident) error {
	model := &models.Incident{}
	model.FromDomain(incident)

	result := r.db.WithContext(ctx).
		Model(&models.Incident{}).
		Where("id = ?", incident.ID).
		Updates(model)

	if result.Error != nil {
		return fmt.Errorf("failed to update incident: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("incident not found: %d", incident.ID)
	}

	return nil
}

// Resolve marks an incident as resolved
func (r *IncidentRepository) Resolve(ctx context.Context, id int64) error {
	// First get the incident
	var model models.Incident
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("incident not found: %d", id)
		}
		return fmt.Errorf("failed to get incident: %w", err)
	}

	// Convert to domain and resolve
	incident := model.ToDomain()
	incident.Resolve()

	// Convert back and update
	model.FromDomain(incident)

	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("failed to resolve incident: %w", err)
	}

	return nil
}

// DeleteResolvedOlderThan deletes resolved incidents that were resolved before the given time
func (r *IncidentRepository) DeleteResolvedOlderThan(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("status = ? AND resolved_at IS NOT NULL AND resolved_at < ?", domain.IncidentStatusResolved, before).
		Delete(&models.Incident{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old incidents: %w", result.Error)
	}

	return result.RowsAffected, nil
}
