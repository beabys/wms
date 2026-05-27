package notificationgrpc

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	notificationv1 "github.com/beabys/wms/proto/gen/go/notification/v1"

	"github.com/beabys/wms/notification-service/internal/app/ports"
	"github.com/beabys/wms/notification-service/internal/application/notification/usecase"
	"github.com/beabys/wms/pkg/authinterceptor"
)

// withExemptMethods wraps an auth interceptor to skip methods listed as exempt.
func withExemptMethods(next grpc.UnaryServerInterceptor, exemptMethods []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		for _, m := range exemptMethods {
			if info.FullMethod == m {
				return handler(ctx, req)
			}
		}
		return next(ctx, req, info, handler)
	}
}

// GRPCServer wraps the gRPC server with its dependencies.
type GRPCServer struct {
	Config   *ports.Config
	Logger   *zap.Logger
	Service  usecase.NotificationServiceHandler
	Server   *grpc.Server
	Listener net.Listener
}

// NewGRPCServer creates a new GRPCServer.
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

// SetConfig sets the server config.
func (s *GRPCServer) SetConfig(cfg *ports.Config) *GRPCServer {
	s.Config = cfg
	return s
}

// SetLogger sets the logger.
func (s *GRPCServer) SetLogger(logger *zap.Logger) *GRPCServer {
	s.Logger = logger
	return s
}

// SetService sets the notification service handler.
func (s *GRPCServer) SetService(svc usecase.NotificationServiceHandler) *GRPCServer {
	s.Service = svc
	return s
}

// Run starts the gRPC server and handles graceful shutdown.
func (s *GRPCServer) Run(ctx context.Context, wg *errgroup.Group) {
	address := fmt.Sprintf(":%d", s.Config.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.Logger.Fatal("failed to listen", zap.Error(err))
	}

	// Auth interceptor with exempt methods (e.g., health check)
	authInterceptor := authinterceptor.ServiceAuthInterceptor([]string{s.Config.Auth.APIKey})
	wrappedInterceptor := withExemptMethods(authInterceptor, s.Config.Auth.ExemptMethods)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(wrappedInterceptor),
	)

	// Register notification service
	notificationSrv := NewNotificationServer(s.Service)
	notificationv1.RegisterNotificationServiceServer(grpcServer, notificationSrv)

	// Register health service
	grpc_health_v1.RegisterHealthServer(grpcServer, &HealthChecker{})

	// Enable reflection
	reflection.Register(grpcServer)

	s.Server = grpcServer
	s.Listener = listener

	wg.Go(func() error {
		s.Logger.Info("gRPC server starting", zap.String("address", address))
		if err := grpcServer.Serve(listener); err != nil {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		s.Logger.Info("gRPC server shutting down")
		grpcServer.GracefulStop()
		return nil
	})
}
