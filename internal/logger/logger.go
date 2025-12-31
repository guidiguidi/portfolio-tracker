package logger

import (
	"log/slog"
	"os"
)

// New creates a new slog.Logger based on the provided configuration.
func New(level string, isJSON bool) *slog.Logger {
	var log *slog.Logger

	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	if isJSON {
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
		)
	} else {
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
		)
	}

	return log
}
