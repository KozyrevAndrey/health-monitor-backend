package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Load loads configuration from a YAML file
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file
	v.SetConfigFile(configPath)

	// Enable environment variable support
	v.SetEnvPrefix("HEALTH_MONITOR")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Explicitly bind env vars — AutomaticEnv alone doesn't work with Unmarshal
	_ = v.BindEnv("server.host")
	_ = v.BindEnv("server.port")
	_ = v.BindEnv("server.enable_auth")
	_ = v.BindEnv("server.basic_auth_user")
	_ = v.BindEnv("server.basic_auth_pass")
	_ = v.BindEnv("database.path")
	_ = v.BindEnv("database.type")
	_ = v.BindEnv("logging.level")
	_ = v.BindEnv("logging.format")
	_ = v.BindEnv("logging.output")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 15*time.Second)
	v.SetDefault("server.write_timeout", 15*time.Second)
	v.SetDefault("server.shutdown_timeout", 30*time.Second)
	v.SetDefault("server.enable_auth", false)

	// Database defaults
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.path", "./data/health-monitor.db")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", 5*time.Minute)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")

	// Retention defaults
	v.SetDefault("retention.check_results", 30*24*time.Hour) // 30 days
	v.SetDefault("retention.incidents", 90*24*time.Hour)     // 90 days
	v.SetDefault("retention.cleanup_interval", 24*time.Hour) // daily
}

// LoadFromBytes loads configuration from byte slice (useful for testing)
func LoadFromBytes(configData []byte) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetConfigType("yaml")
	if err := v.ReadConfig(strings.NewReader(string(configData))); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}
