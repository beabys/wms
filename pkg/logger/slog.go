package logger

import (
	"log/slog"
	"os"
)

// NewSlogLogger creates a new SlogLogger with the given level.
func NewSlogLogger(level slog.Level) *SlogLogger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return &SlogLogger{log: slog.New(handler)}
}

// GetLogger returns the underlying slog.Logger.
func (l *SlogLogger) GetLogger() any {
	return l.log
}

// Debug logs a debug message.
func (l *SlogLogger) Debug(s string, lf ...LogField) {
	args := logFieldsToSlogArgs(lf)
	l.log.Debug(s, args...)
}

// Info logs an info message.
func (l *SlogLogger) Info(s string, lf ...LogField) {
	args := logFieldsToSlogArgs(lf)
	l.log.Info(s, args...)
}

// Warn logs a warning message.
func (l *SlogLogger) Warn(s string, lf ...LogField) {
	args := logFieldsToSlogArgs(lf)
	l.log.Warn(s, args...)
}

// Error logs an error message.
func (l *SlogLogger) Error(s string, err error, lf ...LogField) {
	args := logFieldsToSlogArgs(lf)
	args = append(args, "error", err)
	l.log.Error(s, args...)
}

// Fatal logs a fatal message and calls os.Exit(1).
func (l *SlogLogger) Fatal(s string, lf ...LogField) {
	args := logFieldsToSlogArgs(lf)
	l.log.Error(s, args...)
	os.Exit(1)
}

func logFieldsToSlogArgs(lfs []LogField) []any {
	args := make([]any, 0, len(lfs)*2)
	for _, lf := range lfs {
		args = append(args, lf.Key, lf.Value)
	}
	return args
}
