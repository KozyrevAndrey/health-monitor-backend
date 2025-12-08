package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"health-monitor/internal/storage"
	"health-monitor/pkg/config"
	"health-monitor/pkg/logger"
)

var (
	// Version information (will be set during build)
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/example.yaml", "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version and exit
	if *showVersion {
		fmt.Printf("Health Monitor\n")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Build Time: %s\n", BuildTime)
		fmt.Printf("  Git Commit: %s\n", GitCommit)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	log := logger.New(logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	})

	log.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("git_commit", GitCommit).
		Msg("Starting Health Monitor")

	// Create context with cancellation for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	// Run application
	if err := run(ctx, cfg, log); err != nil {
		log.Fatal().Err(err).Msg("Application failed")
	}

	log.Info().Msg("Health Monitor stopped")
}

// run contains the main application logic
func run(ctx context.Context, cfg *config.Config, log zerolog.Logger) error {
	log.Info().
		Str("address", cfg.Server.GetAddress()).
		Int("targets", len(cfg.Targets)).
		Int("notifiers", len(cfg.Notifiers)).
		Msg("Configuration loaded")

	// Initialize database
	log.Info().Msg("Initializing database...")
	db, err := storage.New(cfg.Database, log)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database")
		}
	}()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create repositories
	targetRepo := storage.NewTargetRepository(db.DB())
	checkResultRepo := storage.NewCheckResultRepository(db.DB())
	incidentRepo := storage.NewIncidentRepository(db.DB())

	log.Info().Msg("Storage layer initialized")

	// TODO: Initialize remaining components
	// - Checker registry
	// - Scheduler
	// - Alert manager
	// - HTTP server

	// For now, just log that we're ready
	log.Info().
		Interface("target_repo", targetRepo != nil).
		Interface("check_result_repo", checkResultRepo != nil).
		Interface("incident_repo", incidentRepo != nil).
		Msg("All components initialized successfully")

	// Wait for shutdown signal
	<-ctx.Done()

	log.Info().Msg("Shutdown signal received, starting graceful shutdown...")

	// Graceful shutdown sequence:
	// TODO: 1. Stop accepting new HTTP requests
	// TODO: 2. Stop scheduler (no new checks)
	// TODO: 3. Wait for active checks to complete
	// 4. Close database connections (handled by defer above)

	log.Info().Msg("Graceful shutdown completed")

	return nil
}
