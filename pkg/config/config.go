package config

import (
	"fmt"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Retention RetentionConfig `mapstructure:"retention"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	EnableAuth      bool          `mapstructure:"enable_auth"`
	BasicAuthUser   string        `mapstructure:"basic_auth_user"`
	BasicAuthPass   string        `mapstructure:"basic_auth_pass"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type           string `mapstructure:"type"`
	Path           string `mapstructure:"path"`
	MaxOpenConns   int    `mapstructure:"max_open_conns"`
	MaxIdleConns   int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// RetentionConfig represents data retention configuration
type RetentionConfig struct {
	CheckResults time.Duration `mapstructure:"check_results"`
	Incidents    time.Duration `mapstructure:"incidents"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Type == "" {
		return fmt.Errorf("database type is required")
	}

	if c.Database.Type == "sqlite" && c.Database.Path == "" {
		return fmt.Errorf("database path is required for SQLite")
	}

	return nil
}

// GetAddress returns the server address in host:port format
func (s *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
