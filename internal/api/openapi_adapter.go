package api

import (
	"net/http"
	"time"

	"health-monitor/internal/domain"
	"health-monitor/internal/generated"
)

// OpenAPIAdapter adapts our existing Server to implement generated.ServerInterface
type OpenAPIAdapter struct {
	server *Server
}

// NewOpenAPIAdapter creates a new adapter
func NewOpenAPIAdapter(server *Server) *OpenAPIAdapter {
	return &OpenAPIAdapter{server: server}
}

// Ensure OpenAPIAdapter implements generated.ServerInterface
var _ generated.ServerInterface = (*OpenAPIAdapter)(nil)

// Health endpoint
func (a *OpenAPIAdapter) GetHealth(w http.ResponseWriter, r *http.Request) {
	a.server.handleHealth(w, r)
}

// Targets endpoints
func (a *OpenAPIAdapter) ListTargets(w http.ResponseWriter, r *http.Request) {
	a.server.handleListTargets(w, r)
}

func (a *OpenAPIAdapter) CreateTarget(w http.ResponseWriter, r *http.Request) {
	a.server.handleCreateTarget(w, r)
}

func (a *OpenAPIAdapter) GetTarget(w http.ResponseWriter, r *http.Request, id string) {
	a.server.handleGetTarget(w, r)
}

func (a *OpenAPIAdapter) UpdateTarget(w http.ResponseWriter, r *http.Request, id string) {
	a.server.handleUpdateTarget(w, r)
}

func (a *OpenAPIAdapter) DeleteTarget(w http.ResponseWriter, r *http.Request, id string) {
	a.server.handleDeleteTarget(w, r)
}

func (a *OpenAPIAdapter) GetTargetResults(w http.ResponseWriter, r *http.Request, id string, params generated.GetTargetResultsParams) {
	a.server.handleGetTargetResults(w, r)
}

func (a *OpenAPIAdapter) GetTargetStats(w http.ResponseWriter, r *http.Request, id string, params generated.GetTargetStatsParams) {
	a.server.handleGetTargetStats(w, r)
}

// Incidents endpoints
func (a *OpenAPIAdapter) ListIncidents(w http.ResponseWriter, r *http.Request, params generated.ListIncidentsParams) {
	a.server.handleListIncidents(w, r)
}

func (a *OpenAPIAdapter) GetIncident(w http.ResponseWriter, r *http.Request, id int64) {
	a.server.handleGetIncident(w, r)
}

func (a *OpenAPIAdapter) ListOngoingIncidents(w http.ResponseWriter, r *http.Request) {
	a.server.handleListOngoingIncidents(w, r)
}

// Helper: Convert generated.Target to domain.Target
func convertToDomainTarget(gt generated.Target) domain.Target {
	var config map[string]interface{}
	if gt.Config != nil {
		config = *gt.Config
	}

	var tags []string
	if gt.Tags != nil {
		tags = *gt.Tags
	}

	var description string
	if gt.Description != nil {
		description = *gt.Description
	}

	return domain.Target{
		ID:          gt.Id,
		Name:        gt.Name,
		Type:        domain.TargetType(gt.Type),
		Config:      config,
		Interval:    time.Duration(gt.Interval),
		Timeout:     time.Duration(gt.Timeout),
		Enabled:     gt.Enabled,
		Tags:        tags,
		Description: description,
		CreatedAt:   safeTime(gt.CreatedAt),
		UpdatedAt:   safeTime(gt.UpdatedAt),
	}
}

// Helper: Convert domain.Target to generated.Target
func convertFromDomainTarget(dt domain.Target) generated.Target {
	tags := dt.Tags
	desc := dt.Description
	config := dt.Config
	createdAt := dt.CreatedAt
	updatedAt := dt.UpdatedAt

	return generated.Target{
		Id:          dt.ID,
		Name:        dt.Name,
		Type:        generated.TargetType(dt.Type),
		Config:      &config,
		Interval:    int64(dt.Interval),
		Timeout:     int64(dt.Timeout),
		Enabled:     dt.Enabled,
		Tags:        &tags,
		Description: &desc,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}
}

// Helper: Safely get time.Time from pointer
func safeTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
