package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"health-monitor/internal/alerting"
	"health-monitor/internal/api"
	"health-monitor/internal/checker"
	"health-monitor/internal/domain"
	"health-monitor/internal/notifier"
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

	alertManager := alerting.NewManager(targetRepo, checkResultRepo, incidentRepo, log)
	log.Info().Msg("Alert manager initialized")

	if err := loadNotifiers(cfg, alertManager, log); err != nil {
		return fmt.Errorf("failed to load notifiers: %w", err)
	}

	sched := scheduler.New(checkerRegistry, checkResultRepo, alertManager, log)

	if err := sched.Start(ctx); err != nil {
		return fmt.Errorf("failed to start scheduler: %w", err)
	}
	defer func() {
		if err := sched.Stop(); err != nil {
			log.Error().Err(err).Msg("Failed to stop scheduler")
		}
	}()

	targets, err := targetRepo.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load targets from database")
	} else {
		for _, target := range targets {
			if !target.Enabled {
				log.Debug().Str("target_id", target.ID).Msg("Skipping disabled target")
				continue
			}

			if err := sched.AddTarget(target); err != nil {
				log.Error().
					Err(err).
					Str("target_id", target.ID).
					Msg("Failed to add target to scheduler")
			}
		}
		log.Info().Int("targets", len(targets)).Msg("Loaded targets from database")
	}

	apiServer := api.NewServer(cfg.Server, targetRepo, checkResultRepo, incidentRepo, sched, log)

	go func() {
		if err := apiServer.Start(); err != nil {
			log.Error().Err(err).Msg("API server failed")
		}
	}()

	log.Info().
		Bool("target_repo", targetRepo != nil).
		Bool("check_result_repo", checkResultRepo != nil).
		Bool("incident_repo", incidentRepo != nil).
		Bool("checker_registry", checkerRegistry != nil).
		Bool("alert_manager", alertManager != nil).
		Bool("scheduler", sched.IsRunning()).
		Bool("api_server", apiServer != nil).
		Msg("All components initialized successfully")

	<-ctx.Done()

	log.Info().Msg("Shutdown signal received, starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("API server shutdown failed")
	}

	log.Info().Msg("Graceful shutdown completed")

	return nil
}

func loadNotifiers(cfg *config.Config, alertManager domain.AlertManager, log zerolog.Logger) error {
	if len(cfg.Notifiers) == 0 {
		log.Info().Msg("No notifiers configured")
		return nil
	}

	enabledCount := 0
	for _, notifierCfg := range cfg.Notifiers {
		if !notifierCfg.Enabled {
			log.Debug().
				Str("id", notifierCfg.ID).
				Str("type", notifierCfg.Type).
				Msg("Notifier disabled, skipping")
			continue
		}

		var n domain.Notifier
		var err error

		switch notifierCfg.Type {
		case "telegram":
			n, err = notifier.NewTelegramNotifier(notifierCfg.Config, log)
		case "email":
			n, err = notifier.NewEmailNotifier(notifierCfg.Config, log)
		default:
			log.Warn().
				Str("id", notifierCfg.ID).
				Str("type", notifierCfg.Type).
				Msg("Unknown notifier type, skipping")
			continue
		}

		if err != nil {
			return fmt.Errorf("failed to create notifier %s (%s): %w", notifierCfg.ID, notifierCfg.Type, err)
		}

		if err := n.Validate(notifierCfg.Config); err != nil {
			return fmt.Errorf("invalid config for notifier %s (%s): %w", notifierCfg.ID, notifierCfg.Type, err)
		}

		alertManager.RegisterNotifier(n)
		enabledCount++

		log.Info().
			Str("id", notifierCfg.ID).
			Str("type", notifierCfg.Type).
			Msg("Notifier registered")
	}

	log.Info().
		Int("enabled", enabledCount).
		Int("total", len(cfg.Notifiers)).
		Msg("Notifiers loaded")

	return nil
}
