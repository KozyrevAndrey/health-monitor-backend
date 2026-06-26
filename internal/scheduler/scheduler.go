package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
	"health-monitor/internal/events"
)

type Scheduler struct {
	checkerRegistry domain.CheckerRegistry
	checkResultRepo domain.CheckResultRepository
	alertManager    domain.AlertManager
	publisher       events.Publisher
	log             zerolog.Logger

	mu      sync.RWMutex
	tasks   map[string]*task
	running bool
	wg      sync.WaitGroup
}

// SetEventPublisher attaches an event publisher so completed checks are
// broadcast in real time. Safe to leave unset (publishing becomes a no-op).
func (s *Scheduler) SetEventPublisher(p events.Publisher) {
	s.publisher = p
}

type task struct {
	target *domain.Target
	ticker *time.Ticker
	cancel context.CancelFunc
}

func New(
	checkerRegistry domain.CheckerRegistry,
	checkResultRepo domain.CheckResultRepository,
	alertManager domain.AlertManager,
	log zerolog.Logger,
) *Scheduler {
	return &Scheduler{
		checkerRegistry: checkerRegistry,
		checkResultRepo: checkResultRepo,
		alertManager:    alertManager,
		log:             log,
		tasks:           make(map[string]*task),
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("scheduler already running")
	}
	s.running = true
	s.mu.Unlock()

	s.log.Info().Msg("Scheduler started")
	return nil
}

func (s *Scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("scheduler not running")
	}

	s.log.Info().Msg("Stopping scheduler...")

	for targetID, t := range s.tasks {
		s.log.Debug().Str("target_id", targetID).Msg("Stopping task")
		t.ticker.Stop()
		if t.cancel != nil {
			t.cancel()
		}
	}

	s.wg.Wait()
	s.running = false

	s.log.Info().Msg("Scheduler stopped")
	return nil
}

func (s *Scheduler) AddTarget(target *domain.Target) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[target.ID]; exists {
		return fmt.Errorf("target %s already scheduled", target.ID)
	}

	if !target.Enabled {
		s.log.Debug().
			Str("target_id", target.ID).
			Msg("Target disabled, skipping")
		return nil
	}

	checker, err := s.checkerRegistry.Get(target.Type)
	if err != nil {
		return fmt.Errorf("failed to get checker for target %s: %w", target.ID, err)
	}

	if err := checker.Validate(target.Config); err != nil {
		return fmt.Errorf("invalid config for target %s: %w", target.ID, err)
	}

	taskCtx, cancel := context.WithCancel(context.Background())

	t := &task{
		target: target,
		ticker: time.NewTicker(target.Interval),
		cancel: cancel,
	}

	s.tasks[target.ID] = t

	s.wg.Add(1)
	go s.runTask(taskCtx, t, checker)

	s.log.Info().
		Str("target_id", target.ID).
		Str("name", target.Name).
		Str("type", string(target.Type)).
		Dur("interval", target.Interval).
		Msg("Target scheduled")

	return nil
}

func (s *Scheduler) RemoveTarget(targetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, exists := s.tasks[targetID]
	if !exists {
		return fmt.Errorf("target %s not found", targetID)
	}

	t.ticker.Stop()
	if t.cancel != nil {
		t.cancel()
	}

	delete(s.tasks, targetID)

	s.log.Info().
		Str("target_id", targetID).
		Msg("Target removed from scheduler")

	return nil
}

func (s *Scheduler) UpdateTarget(target *domain.Target) error {
	if err := s.RemoveTarget(target.ID); err != nil {
		return err
	}
	return s.AddTarget(target)
}

func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *Scheduler) runTask(ctx context.Context, t *task, checker domain.Checker) {
	defer s.wg.Done()

	s.performCheck(ctx, t.target, checker)

	for {
		select {
		case <-ctx.Done():
			s.log.Debug().
				Str("target_id", t.target.ID).
				Msg("Task context cancelled")
			return
		case <-t.ticker.C:
			s.performCheck(ctx, t.target, checker)
		}
	}
}

func (s *Scheduler) performCheck(ctx context.Context, target *domain.Target, checker domain.Checker) {
	checkCtx, cancel := context.WithTimeout(ctx, target.Timeout+5*time.Second)
	defer cancel()

	s.log.Debug().
		Str("target_id", target.ID).
		Str("name", target.Name).
		Msg("Performing check")

	result, err := checker.Check(checkCtx, target)
	if err != nil {
		s.log.Error().
			Err(err).
			Str("target_id", target.ID).
			Msg("Check failed")
		return
	}

	if err := s.checkResultRepo.Save(checkCtx, result); err != nil {
		s.log.Error().
			Err(err).
			Str("target_id", target.ID).
			Msg("Failed to save check result")
		return
	}

	if s.alertManager != nil {
		if err := s.alertManager.ProcessCheckResult(checkCtx, result); err != nil {
			s.log.Error().
				Err(err).
				Str("target_id", target.ID).
				Msg("Failed to process check result in alert manager")
		}
	}

	if s.publisher != nil {
		s.publisher.Publish(events.Event{
			Type: "check",
			Data: map[string]interface{}{
				"target_id":        result.TargetID,
				"target_name":      target.Name,
				"status":           result.Status,
				"response_time_ms": result.ResponseTimeMs,
				"message":          result.Message,
				"checked_at":       result.CheckedAt,
			},
		})
	}

	s.log.Info().
		Str("target_id", target.ID).
		Str("name", target.Name).
		Str("status", string(result.Status)).
		Int64("response_time_ms", result.ResponseTimeMs).
		Msg("Check completed")
}
