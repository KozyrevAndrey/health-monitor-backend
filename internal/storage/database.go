package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"health-monitor/internal/storage/models"
	"health-monitor/pkg/config"
)

// Database manages database connection and operations
type Database struct {
	db  *gorm.DB
	log zerolog.Logger
}

// New creates a new database instance
func New(cfg config.DatabaseConfig, log zerolog.Logger) (*Database, error) {
	var dialector gorm.Dialector

	switch cfg.Type {
	case "sqlite":
		// Ensure directory exists for SQLite file
		if err := ensureDir(cfg.Path); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
		dialector = sqlite.Open(cfg.Path)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// Configure GORM logger
	gormLogger := logger.New(
		&gormLogWriter{log: log},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Get underlying SQL DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure connection pool
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	log.Info().
		Str("type", cfg.Type).
		Str("path", cfg.Path).
		Msg("Database connection established")

	return &Database{
		db:  db,
		log: log,
	}, nil
}

// AutoMigrate runs automatic migrations
func (d *Database) AutoMigrate() error {
	d.log.Info().Msg("Running database migrations...")

	if err := d.db.AutoMigrate(
		&models.Target{},
		&models.CheckResult{},
		&models.Incident{},
	); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	d.log.Info().Msg("Database migrations completed")
	return nil
}

// DB returns the underlying GORM database instance
func (d *Database) DB() *gorm.DB {
	return d.db
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	d.log.Info().Msg("Database connection closed")
	return nil
}

// Ping checks if database is reachable
func (d *Database) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// gormLogWriter adapts zerolog to GORM logger interface
type gormLogWriter struct {
	log zerolog.Logger
}

func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	w.log.Debug().Msgf(format, args...)
}

// ensureDir ensures the directory for a file path exists
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
