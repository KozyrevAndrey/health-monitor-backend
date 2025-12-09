package scheduler

import (
	"fmt"
	"time"

	"health-monitor/internal/domain"
	"health-monitor/pkg/config"
)

func LoadTargetsFromConfig(cfgTargets []config.TargetConfig) ([]*domain.Target, error) {
	targets := make([]*domain.Target, 0, len(cfgTargets))

	for _, cfgTarget := range cfgTargets {
		target, err := configTargetToDomain(cfgTarget)
		if err != nil {
			return nil, fmt.Errorf("failed to convert target %s: %w", cfgTarget.ID, err)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func configTargetToDomain(cfg config.TargetConfig) (*domain.Target, error) {
	interval, err := time.ParseDuration(cfg.Interval)
	if err != nil {
		return nil, fmt.Errorf("invalid interval %s: %w", cfg.Interval, err)
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout %s: %w", cfg.Timeout, err)
	}

	target := &domain.Target{
		ID:          cfg.ID,
		Name:        cfg.Name,
		Type:        domain.TargetType(cfg.Type),
		Config:      cfg.Config,
		Interval:    interval,
		Timeout:     timeout,
		Enabled:     cfg.Enabled,
		Tags:        cfg.Tags,
		Description: cfg.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return target, nil
}
