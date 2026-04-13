package storage

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"health-monitor/internal/domain"
	"health-monitor/internal/storage/models"
)

// NotifierRepository implements domain.NotifierRepository
type NotifierRepository struct {
	db *gorm.DB
}

// NewNotifierRepository creates a new notifier repository
func NewNotifierRepository(db *gorm.DB) *NotifierRepository {
	return &NotifierRepository{db: db}
}

// Create creates a new notifier config
func (r *NotifierRepository) Create(ctx context.Context, cfg *domain.NotifierConfig) error {
	model := &models.NotifierConfig{}
	if err := model.FromDomain(cfg); err != nil {
		return fmt.Errorf("failed to convert notifier config to model: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create notifier config: %w", err)
	}

	return nil
}

// Get retrieves a notifier config by ID
func (r *NotifierRepository) Get(ctx context.Context, id string) (*domain.NotifierConfig, error) {
	var model models.NotifierConfig
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("notifier config not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get notifier config: %w", err)
	}

	cfg, err := model.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to notifier config: %w", err)
	}

	return cfg, nil
}

// List retrieves all notifier configs
func (r *NotifierRepository) List(ctx context.Context) ([]*domain.NotifierConfig, error) {
	var modelCfgs []models.NotifierConfig
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&modelCfgs).Error; err != nil {
		return nil, fmt.Errorf("failed to list notifier configs: %w", err)
	}

	cfgs := make([]*domain.NotifierConfig, 0, len(modelCfgs))
	for _, model := range modelCfgs {
		cfg, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to notifier config: %w", err)
		}
		cfgs = append(cfgs, cfg)
	}

	return cfgs, nil
}

// Update updates an existing notifier config
func (r *NotifierRepository) Update(ctx context.Context, cfg *domain.NotifierConfig) error {
	model := &models.NotifierConfig{}
	if err := model.FromDomain(cfg); err != nil {
		return fmt.Errorf("failed to convert notifier config to model: %w", err)
	}

	result := r.db.WithContext(ctx).
		Model(&models.NotifierConfig{}).
		Where("id = ?", cfg.ID).
		Updates(map[string]interface{}{
			"type":    model.Type,
			"enabled": model.Enabled,
			"config":  model.Config,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update notifier config: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notifier config not found: %s", cfg.ID)
	}

	return nil
}

// Delete deletes a notifier config by ID
func (r *NotifierRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.NotifierConfig{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete notifier config: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notifier config not found: %s", id)
	}

	return nil
}
