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
	log             zerolog.Logger
	enableAuth      bool
	basicAuthUser   string
	basicAuthPass   string
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
		basicAuthUser:   cfg.BasicAuthUser,
		basicAuthPass:   cfg.BasicAuthPass,
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

func (s *Server) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != s.basicAuthUser || pass != s.basicAuthPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Health Monitor"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(s.loggingMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	if s.enableAuth {
		r.Use(s.basicAuthMiddleware)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Web UI and static files
	r.Get("/", s.handleIndex)
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Swagger UI - serves interactive API documentation
	r.Get("/swagger", s.handleSwaggerUI)
	r.Get("/openapi.yaml", s.handleOpenAPISpec)

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
