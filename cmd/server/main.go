package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"health-monitor/internal/checker"
	"health-monitor/internal/scheduler"
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

	checkerRegistry := checker.NewDefaultRegistry()
	log.Info().
		Int("checkers", len(checkerRegistry.List())).
		Interface("types", checkerRegistry.List()).
		Msg("Checker registry initialized")

	sched := scheduler.New(checkerRegistry, checkResultRepo, log)

	targets, err := scheduler.LoadTargetsFromConfig(cfg.Targets)
	if err != nil {
		return fmt.Errorf("failed to load targets from config: %w", err)
	}

	log.Info().Int("targets", len(targets)).Msg("Loaded targets from configuration")

	if err := sched.Start(ctx); err != nil {
		return fmt.Errorf("failed to start scheduler: %w", err)
	}
	defer func() {
		if err := sched.Stop(); err != nil {
			log.Error().Err(err).Msg("Failed to stop scheduler")
		}
	}()

	for _, target := range targets {
		if err := targetRepo.Create(ctx, target); err != nil {
			log.Warn().
				Err(err).
				Str("target_id", target.ID).
				Msg("Failed to save target to database (might already exist)")
		}

		if err := sched.AddTarget(target); err != nil {
			log.Error().
				Err(err).
				Str("target_id", target.ID).
				Msg("Failed to add target to scheduler")
		}
	}

	log.Info().
		Bool("target_repo", targetRepo != nil).
		Bool("check_result_repo", checkResultRepo != nil).
		Bool("incident_repo", incidentRepo != nil).
		Bool("checker_registry", checkerRegistry != nil).
		Bool("scheduler", sched.IsRunning()).
		Msg("All components initialized successfully")

	<-ctx.Done()

	log.Info().Msg("Shutdown signal received, starting graceful shutdown...")

	log.Info().Msg("Graceful shutdown completed")

	return nil
}
