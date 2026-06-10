package grpcdapter

import (
	"context"
	"log/slog"
	"testing"

	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestLoggingInterceptor_Success(t *testing.T) {
	log := logger.NewSlogLogger(slog.LevelDebug)
	interceptor := LoggingInterceptor(log)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
		FullMethod: "/test.v1.TestService/Method",
	}, handler)

	require.NoError(t, err)
	assert.Equal(t, "response", resp)
}

func TestLoggingInterceptor_Error(t *testing.T) {
	log := logger.NewSlogLogger(slog.LevelDebug)
	interceptor := LoggingInterceptor(log)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, assert.AnError
	}

	resp, err := interceptor(context.Background(), "request", &grpc.UnaryServerInfo{
		FullMethod: "/test.v1.TestService/FailingMethod",
	}, handler)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, assert.AnError)
}


