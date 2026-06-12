package logger

import (
	"log/slog"

	"go.uber.org/zap"
)

// Logger defines the interface for structured logging.
type Logger interface {
	GetLogger() any
	Debug(string, ...LogField)
	Info(string, ...LogField)
	Warn(string, ...LogField)
	Error(string, error, ...LogField)
	Fatal(string, ...LogField)
}

// ZapLogger implements Logger using zap.
type ZapLogger struct {
	log *zap.Logger
}

// SlogLogger implements Logger using slog.
type SlogLogger struct {
	log *slog.Logger
}

// LogField is a key-value pair for structured logging.
type LogField struct {
	Key   string
	Value any
}
