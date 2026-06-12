package grpcdapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func TestHealthServerCheck(t *testing.T) {
	s := NewHealthServer()
	require.NotNil(t, s)

	resp, err := s.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthServerCheckWithService(t *testing.T) {
	s := NewHealthServer()
	resp, err := s.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{
		Service: "auth",
	})
	require.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthServerWatch(t *testing.T) {
	s := NewHealthServer()
	err := s.Watch(&grpc_health_v1.HealthCheckRequest{}, nil)
	require.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unimplemented, st.Code())
	assert.Contains(t, st.Message(), "watch not implemented")
}
