// Package retention provides a background job that periodically deletes old
// check results and resolved incidents according to the configured retention
// policy, keeping the database from growing without bound.
package retention

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"health-monitor/internal/domain"
	"health-monitor/pkg/config"
)

// Cleaner periodically removes data older than the configured retention windows.
type Cleaner struct {
	checkResults domain.CheckResultRepository
	incidents    domain.IncidentRepository
	cfg          config.RetentionConfig
	log          zerolog.Logger
}

// New creates a retention Cleaner.
func New(
	checkResults domain.CheckResultRepository,
	incidents domain.IncidentRepository,
	cfg config.RetentionConfig,
	log zerolog.Logger,
) *Cleaner {
	return &Cleaner{
		checkResults: checkResults,
		incidents:    incidents,
		cfg:          cfg,
		log:          log.With().Str("component", "retention").Logger(),
	}
}

// Start launches the cleanup loop in a goroutine. It runs an initial sweep
// shortly after start and then every cleanup_interval until ctx is cancelled.
// If cleanup_interval is not positive the cleaner is disabled.
func (c *Cleaner) Start(ctx context.Context) {
	if c.cfg.CleanupInterval <= 0 {
		c.log.Info().Msg("Retention cleanup disabled (cleanup_interval <= 0)")
		return
	}

	c.log.Info().
		Dur("interval", c.cfg.CleanupInterval).
		Dur("check_results_ttl", c.cfg.CheckResults).
		Dur("incidents_ttl", c.cfg.Incidents).
		Msg("Retention cleanup started")

	go c.loop(ctx)
}

func (c *Cleaner) loop(ctx context.Context) {
	ticker := time.NewTicker(c.cfg.CleanupInterval)
	defer ticker.Stop()

	// Run an initial sweep so a long restart cycle does not let data pile up.
	c.cleanup(ctx)

	for {
		select {
		case <-ctx.Done():
			c.log.Info().Msg("Retention cleanup stopped")
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

// cleanup performs a single retention sweep. Each retention window with a
// non-positive duration is skipped (treated as "keep forever").
func (c *Cleaner) cleanup(ctx context.Context) {
	now := time.Now()

	if c.cfg.CheckResults > 0 {
		before := now.Add(-c.cfg.CheckResults)
		if deleted, err := c.checkResults.DeleteOlderThan(ctx, before); err != nil {
			c.log.Error().Err(err).Msg("Failed to delete old check results")
		} else if deleted > 0 {
			c.log.Info().Int64("deleted", deleted).Time("before", before).Msg("Deleted old check results")
		}
	}

	if c.cfg.Incidents > 0 {
		before := now.Add(-c.cfg.Incidents)
		if deleted, err := c.incidents.DeleteResolvedOlderThan(ctx, before); err != nil {
			c.log.Error().Err(err).Msg("Failed to delete old incidents")
		} else if deleted > 0 {
			c.log.Info().Int64("deleted", deleted).Time("before", before).Msg("Deleted old resolved incidents")
		}
	}
}
