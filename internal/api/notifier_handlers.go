package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"health-monitor/internal/domain"
	"health-monitor/internal/notifier"
)

func (s *Server) handleListNotifiers(w http.ResponseWriter, r *http.Request) {
	cfgs, err := s.notifierRepo.List(r.Context())
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to list notifiers", err)
		return
	}

	masked := make([]*domain.NotifierConfig, 0, len(cfgs))
	for _, cfg := range cfgs {
		masked = append(masked, maskNotifierConfig(cfg))
	}

	s.respondJSON(w, http.StatusOK, masked)
}

func (s *Server) handleGetNotifier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cfg, err := s.notifierRepo.Get(r.Context(), id)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Notifier not found", err)
		return
	}

	s.respondJSON(w, http.StatusOK, maskNotifierConfig(cfg))
}

func (s *Server) handleCreateNotifier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var cfg domain.NotifierConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}
	cfg.CreatedAt = time.Now()
	cfg.UpdatedAt = time.Now()

	if _, err := buildNotifier(&cfg, s); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid notifier config", err)
		return
	}

	if err := s.notifierRepo.Create(ctx, &cfg); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to create notifier", err)
		return
	}

	if err := s.reloadNotifiers(ctx); err != nil {
		s.log.Error().Err(err).Msg("Failed to reload notifiers after create")
	}

	s.respondJSON(w, http.StatusCreated, maskNotifierConfig(&cfg))
}

func (s *Server) handleUpdateNotifier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	existing, err := s.notifierRepo.Get(ctx, id)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Notifier not found", err)
		return
	}

	var updates domain.NotifierConfig
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	existing.Type = updates.Type
	existing.Enabled = updates.Enabled
	existing.UpdatedAt = time.Now()

	// Merge config fields — keep existing values for empty/masked fields
	if updates.Config != nil {
		for k, v := range updates.Config {
			strVal, isStr := v.(string)
			if isStr && (strVal == "" || strVal == "***") {
				continue
			}
			existing.Config[k] = v
		}
	}

	if _, err := buildNotifier(existing, s); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid notifier config", err)
		return
	}

	if err := s.notifierRepo.Update(ctx, existing); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to update notifier", err)
		return
	}

	if err := s.reloadNotifiers(ctx); err != nil {
		s.log.Error().Err(err).Msg("Failed to reload notifiers after update")
	}

	s.respondJSON(w, http.StatusOK, maskNotifierConfig(existing))
}

func (s *Server) handleDeleteNotifier(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := s.notifierRepo.Delete(ctx, id); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to delete notifier", err)
		return
	}

	if err := s.reloadNotifiers(ctx); err != nil {
		s.log.Error().Err(err).Msg("Failed to reload notifiers after delete")
	}

	w.WriteHeader(http.StatusNoContent)
}

// reloadNotifiers loads all enabled notifiers from DB and re-registers with AlertManager
func (s *Server) reloadNotifiers(ctx context.Context) error {
	if s.alertManager == nil || s.notifierRepo == nil {
		return nil
	}

	cfgs, err := s.notifierRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list notifier configs: %w", err)
	}

	s.alertManager.ClearNotifiers()

	for _, cfg := range cfgs {
		if !cfg.Enabled {
			continue
		}

		n, err := buildNotifier(cfg, s)
		if err != nil {
			s.log.Error().
				Err(err).
				Str("id", cfg.ID).
				Str("type", cfg.Type).
				Msg("Failed to build notifier, skipping")
			continue
		}

		s.alertManager.RegisterNotifier(n)

		s.log.Info().
			Str("id", cfg.ID).
			Str("type", cfg.Type).
			Msg("Notifier registered")
	}

	return nil
}

// maskNotifierConfig returns a copy with sensitive config fields masked
func maskNotifierConfig(cfg *domain.NotifierConfig) *domain.NotifierConfig {
	masked := &domain.NotifierConfig{
		ID:        cfg.ID,
		Type:      cfg.Type,
		Enabled:   cfg.Enabled,
		CreatedAt: cfg.CreatedAt,
		UpdatedAt: cfg.UpdatedAt,
		Config:    make(map[string]interface{}),
	}

	sensitiveKeys := map[string]bool{
		"smtp_password": true,
		"bot_token":     true,
		"password":      true,
		"token":         true,
		"secret":        true,
	}

	for k, v := range cfg.Config {
		if sensitiveKeys[k] {
			masked.Config[k] = "***"
		} else {
			masked.Config[k] = v
		}
	}

	return masked
}

// buildNotifier creates a domain.Notifier from a NotifierConfig
func buildNotifier(cfg *domain.NotifierConfig, s *Server) (domain.Notifier, error) {
	switch cfg.Type {
	case "email":
		return notifier.NewEmailNotifier(cfg.Config, s.log)
	case "telegram":
		return notifier.NewTelegramNotifier(cfg.Config, s.log)
	case "gmail":
		return notifier.NewGmailNotifier(cfg.Config, s.log)
	case "gmail_oauth":
		return notifier.NewGmailOAuthNotifier(cfg.Config, s.log)
	case "webhook":
		return notifier.NewWebhookNotifier(cfg.Config, s.log)
	default:
		return nil, fmt.Errorf("unknown notifier type: %s", cfg.Type)
	}
}
