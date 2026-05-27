package inventorygrpc

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"

	"github.com/beabys/wms/inventory-service/internal/app/ports"
	"github.com/beabys/wms/inventory-service/internal/application/inventory/usecase"
	"github.com/beabys/wms/pkg/authinterceptor"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	inventoryv1 "github.com/beabys/wms/proto/gen/go/inventory/v1"
)

// GRPCServer wraps the gRPC server with its dependencies.
type GRPCServer struct {
	Config   *ports.Config
	Logger   *zap.Logger
	Service  usecase.InventoryServiceHandler
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

// SetService sets the inventory service handler.
func (s *GRPCServer) SetService(svc usecase.InventoryServiceHandler) *GRPCServer {
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

	// Load public key for JWT validation
	var publicKey *rsa.PublicKey
	if s.Config.JWTPublicKeyPath != "" {
		pemBytes, err := os.ReadFile(s.Config.JWTPublicKeyPath)
		if err != nil {
			s.Logger.Fatal("failed to read public key", zap.Error(err))
		}
		block, _ := pem.Decode(pemBytes)
		if block == nil {
			s.Logger.Fatal("no PEM block in public key")
		}
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			s.Logger.Fatal("failed to parse public key", zap.Error(err))
		}
		var ok bool
		publicKey, ok = key.(*rsa.PublicKey)
		if !ok {
			s.Logger.Fatal("public key is not RSA")
		}
	}

	grpcOptions := []grpc.ServerOption{}
	if publicKey != nil {
		grpcOptions = append(grpcOptions, grpc.UnaryInterceptor(
			authinterceptor.JWTAuthInterceptor(publicKey, s.Config.Auth.ExemptMethods),
		))
	}
	grpcServer := grpc.NewServer(grpcOptions...)

	// Register inventory service
	inventorySrv := NewInventoryServer(s.Service)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, inventorySrv)

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
