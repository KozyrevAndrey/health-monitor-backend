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
	"health-monitor/pkg/config"
)

type Server struct {
	server          *http.Server
	router          chi.Router
	targetRepo      domain.TargetRepository
	checkResultRepo domain.CheckResultRepository
	incidentRepo    domain.IncidentRepository
	log             zerolog.Logger
}

func NewServer(
	cfg config.ServerConfig,
	targetRepo domain.TargetRepository,
	checkResultRepo domain.CheckResultRepository,
	incidentRepo domain.IncidentRepository,
	log zerolog.Logger,
) *Server {
	s := &Server{
		targetRepo:      targetRepo,
		checkResultRepo: checkResultRepo,
		incidentRepo:    incidentRepo,
		log:             log,
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

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(s.loggingMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", s.handleHealth)

	r.Get("/", s.handleIndex)
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/targets", func(r chi.Router) {
			r.Get("/", s.handleListTargets)
			r.Post("/", s.handleCreateTarget)
			r.Get("/{id}", s.handleGetTarget)
			r.Put("/{id}", s.handleUpdateTarget)
			r.Delete("/{id}", s.handleDeleteTarget)
			r.Get("/{id}/results", s.handleGetTargetResults)
			r.Get("/{id}/stats", s.handleGetTargetStats)
		})

		r.Route("/results", func(r chi.Router) {
			r.Get("/", s.handleListResults)
			r.Get("/{id}", s.handleGetResult)
		})

		r.Route("/incidents", func(r chi.Router) {
			r.Get("/", s.handleListIncidents)
			r.Get("/{id}", s.handleGetIncident)
			r.Get("/ongoing", s.handleListOngoingIncidents)
		})
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
