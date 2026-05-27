package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new zap.Logger configured for the given environment.
// In production, it uses JSON encoding; otherwise, console encoding.
func New(env string) (*zap.Logger, error) {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return logger, nil
}

// MustNew creates a new logger and panics on error.
func MustNew(env string) *zap.Logger {
	logger, err := New(env)
	if err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	return logger
}

// Sync flushes any buffered log entries. Call as defer logger.Sync().
func Sync(logger *zap.Logger) {
	_ = logger.Sync()
}

// Ctx extracts a logger from context or returns the fallback.
func Ctx(ctx context.Context, fallback *zap.Logger) *zap.Logger {
	if ctx == nil {
		return fallback
	}
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return fallback
}

type ctxKey struct{}
