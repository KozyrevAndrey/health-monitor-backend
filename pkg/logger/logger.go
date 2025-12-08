package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config represents logger configuration
type Config struct {
	Level  string // debug, info, warn, error
	Format string // json, console
	Output string // stdout, stderr, or file path
}

// New creates a new configured logger
func New(cfg Config) zerolog.Logger {
	// Parse log level
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	// Determine output writer
	output := getOutput(cfg.Output)

	// Configure logger based on format
	var logger zerolog.Logger
	if cfg.Format == "console" {
		// Pretty console output with colors
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		// JSON format (default)
		logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	}

	// Set as global logger
	log.Logger = logger

	return logger
}

// parseLevel converts string level to zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// getOutput returns the appropriate writer based on output string
func getOutput(output string) io.Writer {
	switch strings.ToLower(output) {
	case "stdout", "":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		// Assume it's a file path
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file %s: %v, falling back to stdout\n", output, err)
			return os.Stdout
		}
		return file
	}
}

// Default creates a logger with default settings (useful for quick setup)
func Default() zerolog.Logger {
	return New(Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	})
}

// WithFields creates a logger with additional fields
func WithFields(logger zerolog.Logger, fields map[string]interface{}) zerolog.Logger {
	ctx := logger.With()
	for key, value := range fields {
		ctx = ctx.Interface(key, value)
	}
	return ctx.Logger()
}
