package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
	"health-monitor/internal/events"
	"health-monitor/internal/generated"
	"health-monitor/pkg/config"
)

type Server struct {
	server          *http.Server
	router          chi.Router
	targetRepo      domain.TargetRepository
	checkResultRepo domain.CheckResultRepository
	incidentRepo    domain.IncidentRepository
	notifierRepo    domain.NotifierRepository
	alertManager    domain.AlertManager
	scheduler       domain.Scheduler
	eventBroker     *events.Broker
	dbPinger        pinger
	log             zerolog.Logger
	enableAuth      bool
	authUser        string
	authPass        string
	sessions        *sessionStore
}

// pinger is the minimal dependency the health check needs from the database.
type pinger interface {
	Ping() error
}

// SetEventBroker attaches the broker used by the SSE /api/v1/events endpoint.
func (s *Server) SetEventBroker(b *events.Broker) {
	s.eventBroker = b
}

// SetDBPinger attaches a database health checker used by the /health endpoint.
func (s *Server) SetDBPinger(p pinger) {
	s.dbPinger = p
}

func NewServer(
	cfg config.ServerConfig,
	targetRepo domain.TargetRepository,
	checkResultRepo domain.CheckResultRepository,
	incidentRepo domain.IncidentRepository,
	notifierRepo domain.NotifierRepository,
	alertManager domain.AlertManager,
	scheduler domain.Scheduler,
	log zerolog.Logger,
) *Server {
	s := &Server{
		targetRepo:      targetRepo,
		checkResultRepo: checkResultRepo,
		incidentRepo:    incidentRepo,
		notifierRepo:    notifierRepo,
		alertManager:    alertManager,
		scheduler:       scheduler,
		log:             log,
		enableAuth:      cfg.EnableAuth,
		authUser:        cfg.AuthUser,
		authPass:        cfg.AuthPass,
		sessions:        newSessionStore(),
	}

	s.router = s.setupRouter()

	s.server = &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      s.router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return s
}

func (s *Server) setupRouter() chi.Router {
	r := chi.NewRouter()

	// Common middleware applied to every route. NOTE: neither the logging
	// wrapper nor a request Timeout is here — both break long-lived SSE
	// connections (the logging wrapper can prevent clearing the write
	// deadline; Timeout would close the stream after 60s and mask client
	// disconnect). They are applied to the grouped routes below instead.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	if s.enableAuth {
		r.Use(s.sessionAuthMiddleware)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Server-Sent Events stream — must stay outside the logging/Timeout group.
	r.Get("/api/v1/events", s.handleEvents)

	// All other routes get request logging and a 60s timeout.
	r.Group(func(r chi.Router) {
		r.Use(s.loggingMiddleware)
		r.Use(middleware.Timeout(60 * time.Second))

		// Web UI and static files
		r.Get("/", s.handleIndex)
		fileServer := http.FileServer(http.Dir("./web/static"))
		r.Handle("/static/*", http.StripPrefix("/static", fileServer))

		// Session-based authentication
		r.Get("/login", s.handleLoginPage)
		r.Post("/api/v1/auth/login", s.handleLogin)
		r.Post("/api/v1/auth/logout", s.handleLogout)
		r.Get("/api/v1/auth/me", s.handleAuthStatus)

		// Swagger UI - serves interactive API documentation
		r.Get("/swagger", s.handleSwaggerUI)
		r.Get("/openapi.yaml", s.handleOpenAPISpec)

		// Per-target incidents (not in OpenAPI spec)
		r.Get("/api/v1/targets/{id}/incidents", s.handleGetTargetIncidents)

		// Notifier CRUD routes (not in OpenAPI spec)
		r.Route("/api/v1/notifiers", func(r chi.Router) {
			r.Get("/", s.handleListNotifiers)
			r.Post("/", s.handleCreateNotifier)
			r.Get("/{id}", s.handleGetNotifier)
			r.Put("/{id}", s.handleUpdateNotifier)
			r.Delete("/{id}", s.handleDeleteNotifier)
		})

		// Create OpenAPI adapter and register all API routes automatically
		// This replaces manual route registration with generated routes from openapi.yaml
		adapter := NewOpenAPIAdapter(s)
		generated.HandlerFromMux(adapter, r)
	})

	return r
}

func (s *Server) Start() error {
	s.log.Info().
		Str("address", s.server.Addr).
		Msg("Starting HTTP API server")

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info().Msg("Shutting down HTTP API server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.log.Info().Msg("HTTP API server stopped")
	return nil
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		s.log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.Status()).
			Int("bytes", ww.BytesWritten()).
			Dur("duration_ms", time.Since(start)).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("HTTP request")
	})
}
