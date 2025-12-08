package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"health-monitor/internal/domain"
	"health-monitor/internal/storage/models"
)

// TargetRepository implements domain.TargetRepository
type TargetRepository struct {
	db *gorm.DB
}

// NewTargetRepository creates a new target repository
func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

// Create creates a new target
func (r *TargetRepository) Create(ctx context.Context, target *domain.Target) error {
	model := &models.Target{}
	if err := model.FromDomain(target); err != nil {
		return fmt.Errorf("failed to convert target to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}

	return nil
}

// Get retrieves a target by ID
func (r *TargetRepository) Get(ctx context.Context, id string) (*domain.Target, error) {
	var model models.Target
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("target not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	target, err := model.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to target: %w", err)
	}

	return target, nil
}

// List retrieves all targets
func (r *TargetRepository) List(ctx context.Context) ([]*domain.Target, error) {
	var modelTargets []models.Target
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&modelTargets).Error; err != nil {
		return nil, fmt.Errorf("failed to list targets: %w", err)
	}

	targets := make([]*domain.Target, 0, len(modelTargets))
	for _, model := range modelTargets {
		target, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to target: %w", err)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

// ListEnabled retrieves all enabled targets
func (r *TargetRepository) ListEnabled(ctx context.Context) ([]*domain.Target, error) {
	var modelTargets []models.Target
	if err := r.db.WithContext(ctx).
		Where("enabled = ?", true).
		Order("created_at DESC").
		Find(&modelTargets).Error; err != nil {
		return nil, fmt.Errorf("failed to list enabled targets: %w", err)
	}

	targets := make([]*domain.Target, 0, len(modelTargets))
	for _, model := range modelTargets {
		target, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to target: %w", err)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

// Update updates an existing target
func (r *TargetRepository) Update(ctx context.Context, target *domain.Target) error {
	model := &models.Target{}
	if err := model.FromDomain(target); err != nil {
		return fmt.Errorf("failed to convert target to model: %w", err)
	}

	result := r.db.WithContext(ctx).
		Model(&models.Target{}).
		Where("id = ?", target.ID).
		Updates(model)

	if result.Error != nil {
		return fmt.Errorf("failed to update target: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("target not found: %s", target.ID)
	}

	return nil
}

// Delete deletes a target by ID
func (r *TargetRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Target{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete target: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("target not found: %s", id)
	}

	return nil
}
