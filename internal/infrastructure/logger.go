package infrastructure

import (
	"log/slog"
	"os"
	"strings"

	"github.com/dokkiitech/grumble-back/internal/logging"
)

// NewLogger creates a new slog-based logger.
//
// Configuration via env vars:
// - LOG_FORMAT: "json" for JSON logs (default), "text" for human-readable
// - LOG_LEVEL:  "debug", "info" (default), "warn", "error"
func NewLogger() logging.Logger {
	// Select handler (JSON by default)
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))

	// Select level (info by default)
	var level slog.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level, AddSource: true}

	var handler slog.Handler
	switch format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default: // json
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
