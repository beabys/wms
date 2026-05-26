package inboundgrpc

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthChecker implements grpc.health.v1.Health.
type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

// Check returns SERVING always.
func (s *HealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch streams SERVING status.
func (s *HealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
