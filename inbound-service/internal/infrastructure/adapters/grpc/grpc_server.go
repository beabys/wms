package inboundgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/beabys/wms/inbound-service/internal/app/ports"
	"github.com/beabys/wms/inbound-service/internal/application/inbound/usecase"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"
)

// GRPCServer wraps the gRPC server with its dependencies.
type GRPCServer struct {
	Config   *ports.Config
	Logger   *zap.Logger
	Service  usecase.InboundServiceHandler
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

// SetService sets the inbound service handler.
func (s *GRPCServer) SetService(svc usecase.InboundServiceHandler) *GRPCServer {
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

	grpcServer := grpc.NewServer()

	// Register inbound service
	inboundSrv := NewInboundServer(s.Service)
	inboundv1.RegisterInboundServiceServer(grpcServer, inboundSrv)

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
