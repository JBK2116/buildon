// Package logging provides logging related functionality for the application.
package logging

import (
	"log/slog"
	"os"
)

// newLogger returns a new [slog.Logger] that writes JSON formatted logs to standard output.
func newLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}
