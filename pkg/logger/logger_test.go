package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beabys/wms/pkg/logger"
)

func TestNew_Development(t *testing.T) {
	l, err := logger.New("development")
	require.NoError(t, err)
	require.NotNil(t, l)
	l.Info("test dev log")
	logger.Sync(l)
}

func TestNew_Production(t *testing.T) {
	l, err := logger.New("production")
	require.NoError(t, err)
	require.NotNil(t, l)
	l.Info("test prod log")
	logger.Sync(l)
}

func TestNew_InvalidEnv(t *testing.T) {
	l, err := logger.New("invalid")
	require.NoError(t, err)
	require.NotNil(t, l)
	// Falls back to development config defaults
	l.Info("test fallback")
	logger.Sync(l)
}

func TestMustNew(t *testing.T) {
	l := logger.MustNew("development")
	require.NotNil(t, l)
	logger.Sync(l)
}

func TestCtx_WithLogger(t *testing.T) {
	base := logger.MustNew("development")
	ctx := context.Background()

	// No logger in context - returns base
	l := logger.Ctx(ctx, base)
	assert.Equal(t, base, l)
}

func TestCtx_NilContext(t *testing.T) {
	base := logger.MustNew("development")
	l := logger.Ctx(nil, base)
	assert.Equal(t, base, l)
}

func TestCtx_WithCustomLogger(t *testing.T) {
	custom := logger.MustNew("development")

	// Context without a stored logger returns the passed fallback
	ctx := context.Background()
	l := logger.Ctx(ctx, custom)
	assert.Equal(t, custom, l)
}

func TestInterceptorLogger(t *testing.T) {
	zapLogger := logger.MustNew("development")
	il := logger.InterceptorLogger(zapLogger)
	assert.NotNil(t, il)

	// Verify the adapter is functional by logging at various levels
	ctx := context.Background()
	il.Log(ctx, 0, "test debug msg")
	il.Log(ctx, 1, "test info msg")
	il.Log(ctx, 2, "test warn msg")
	il.Log(ctx, 3, "test error msg")
	// Unknown level should not panic
	il.Log(ctx, 99, "test unknown level msg")
	logger.Sync(zapLogger)
}
