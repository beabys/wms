package grpcdapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthServer_Check(t *testing.T) {
	s := NewHealthServer()
	resp, err := s.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.GetStatus())
}

func TestHealthServer_Watch(t *testing.T) {
	s := NewHealthServer()
	err := s.Watch(&grpc_health_v1.HealthCheckRequest{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}
