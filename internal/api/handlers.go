package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"health-monitor/internal/domain"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "dev",
	}
	s.respondJSON(w, http.StatusOK, response)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/index.html")
}

func (s *Server) handleListTargets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targets, err := s.targetRepo.List(ctx)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to list targets", err)
		return
	}

	s.respondJSON(w, http.StatusOK, targets)
}

func (s *Server) handleGetTarget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	target, err := s.targetRepo.Get(ctx, id)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Target not found", err)
		return
	}

	s.respondJSON(w, http.StatusOK, target)
}

func (s *Server) handleCreateTarget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var target domain.Target
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	target.CreatedAt = time.Now()
	target.UpdatedAt = time.Now()

	if err := s.targetRepo.Create(ctx, &target); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to create target", err)
		return
	}

	if target.Enabled && s.scheduler != nil {
		if err := s.scheduler.AddTarget(&target); err != nil {
			s.log.Error().
				Err(err).
				Str("target_id", target.ID).
				Msg("Failed to add target to scheduler")
		}
	}

	s.respondJSON(w, http.StatusCreated, target)
}

func (s *Server) handleUpdateTarget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	existing, err := s.targetRepo.Get(ctx, id)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Target not found", err)
		return
	}

	var updates domain.Target
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	wasEnabled := existing.Enabled

	existing.Name = updates.Name
	existing.Type = updates.Type
	existing.Config = updates.Config
	existing.Interval = updates.Interval
	existing.Timeout = updates.Timeout
	existing.Enabled = updates.Enabled
	existing.Tags = updates.Tags
	existing.Description = updates.Description
	existing.UpdatedAt = time.Now()

	if err := s.targetRepo.Update(ctx, existing); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to update target", err)
		return
	}

	if s.scheduler != nil {
		if wasEnabled && !existing.Enabled {
			s.scheduler.RemoveTarget(existing.ID)
		} else if !wasEnabled && existing.Enabled {
			if err := s.scheduler.AddTarget(existing); err != nil {
				s.log.Error().
					Err(err).
					Str("target_id", existing.ID).
					Msg("Failed to add target to scheduler")
			}
		} else if existing.Enabled {
			if err := s.scheduler.UpdateTarget(existing); err != nil {
				s.log.Error().
					Err(err).
					Str("target_id", existing.ID).
					Msg("Failed to update target in scheduler")
			}
		}
	}

	s.respondJSON(w, http.StatusOK, existing)
}

func (s *Server) handleDeleteTarget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if s.scheduler != nil {
		s.scheduler.RemoveTarget(id)
	}

	if err := s.targetRepo.Delete(ctx, id); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to delete target", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetTargetResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	targetID := chi.URLParam(r, "id")

	limit := 100
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	results, err := s.checkResultRepo.GetHistory(ctx, targetID, limit, offset)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to get results", err)
		return
	}

	s.respondJSON(w, http.StatusOK, results)
}

func (s *Server) handleGetTargetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	targetID := chi.URLParam(r, "id")

	period := 24 * time.Hour
	if p := r.URL.Query().Get("period"); p != "" {
		// NOTE: Go's time.ParseDuration does not understand "d"; callers must
		// use hours (e.g. 168h for 7d, 720h for 30d). Reject bad input rather
		// than silently falling back to the default and returning wrong data.
		parsed, err := time.ParseDuration(p)
		if err != nil || parsed <= 0 {
			s.respondError(w, http.StatusBadRequest, "Invalid period (use Go duration like 24h, 168h, 720h)", err)
			return
		}
		period = parsed
	}

	stats, err := s.checkResultRepo.GetStats(ctx, targetID, period)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to get stats", err)
		return
	}

	s.respondJSON(w, http.StatusOK, stats)
}

// handleGetTargetIncidents returns incidents for a single target, most recent
// first. Registered as a manual route (not part of the OpenAPI spec).
func (s *Server) handleGetTargetIncidents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	targetID := chi.URLParam(r, "id")

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	incidents, err := s.incidentRepo.ListByTarget(ctx, targetID, limit, offset)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to get incidents", err)
		return
	}

	s.respondJSON(w, http.StatusOK, incidents)
}

func (s *Server) handleListResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "List all results - not implemented yet",
	})
}

func (s *Server) handleGetResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Get result by ID - not implemented yet",
	})
}

func (s *Server) handleListIncidents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 100
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	incidents, err := s.incidentRepo.List(ctx, limit, offset)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to list incidents", err)
		return
	}

	s.respondJSON(w, http.StatusOK, incidents)
}

func (s *Server) handleGetIncident(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid incident ID", err)
		return
	}

	incident, err := s.incidentRepo.Get(ctx, id)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Incident not found", err)
		return
	}

	s.respondJSON(w, http.StatusOK, incident)
}

func (s *Server) handleListOngoingIncidents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 100
	offset := 0

	incidents, err := s.incidentRepo.List(ctx, limit, offset)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to list incidents", err)
		return
	}

	ongoing := make([]*domain.Incident, 0)
	for _, incident := range incidents {
		if incident.IsOngoing() {
			ongoing = append(ongoing, incident)
		}
	}

	s.respondJSON(w, http.StatusOK, ongoing)
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.log.Error().Err(err).Msg("Failed to encode JSON response")
		}
	}
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string, err error) {
	s.log.Error().
		Err(err).
		Int("status", status).
		Str("message", message).
		Msg("API error")

	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	s.respondJSON(w, status, response)
}
