package app

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"runtime/debug"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	customerv1 "github.com/beabys/wms/proto/gen/go/customer/v1"

	"github.com/beabys/wms/pkg/authinterceptor"
	"github.com/beabys/wms/pkg/logger"
	"github.com/beabys/wms/pkg/postgres"

	"github.com/beabys/wms/customer-service/internal/app/ports"
	"github.com/beabys/wms/customer-service/internal/application/customer/usecase"
	grpcadapter "github.com/beabys/wms/customer-service/internal/infrastructure/adapters/grpc"
	persrepo "github.com/beabys/wms/customer-service/internal/infrastructure/persistence/repository"
)

// loadPublicKey reads and parses an RSA public key from a PEM file.
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM data in public key file")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA public key")
	}
	return rsaKey, nil
}

// New creates a new empty App.
func New() *App {
	return &App{}
}

// SetConfig loads and sets the application configuration.
func (app *App) SetConfig(config ports.AppConfig) error {
	app.Config = config
	return config.LoadConfig("customer-service/config.yaml")
}

// Setup initializes the app: config, logger, and database pool.
func (app *App) Setup(config ports.AppConfig) error {
	err := app.SetConfig(config)
	if err != nil {
		return err
	}

	log, err := logger.New("development")
	if err != nil {
		return err
	}
	app.Logger = log

	return nil
}

// ConnectDB creates the database connection pool.
func (app *App) ConnectDB(ctx context.Context) error {
	dsn := app.Config.GetDBDSN()
	if dsn == "" {
		return fmt.Errorf("database DSN is empty")
	}

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	app.Pool = pool
	return nil
}

// SetGRPCServer creates and configures the gRPC server.
func (app *App) SetGRPCServer() error {
	customerRepo := persrepo.NewCustomerRepo(app.Pool)
	inviteRepo := persrepo.NewInviteRepo(app.Pool)
	customerService := usecase.NewCustomerService(customerRepo, inviteRepo)

	srv := grpcadapter.NewServer(app.Logger, customerService)

	// Load public key for JWT validation
	pubKeyPath := app.Config.GetJWTPublicKeyPath()
	var publicKey *rsa.PublicKey
	if pubKeyPath != "" {
		var err error
		publicKey, err = loadPublicKey(pubKeyPath)
		if err != nil {
			return fmt.Errorf("load public key: %w", err)
		}
	}

	exemptMethods := []string{
		customerv1.CustomerService_InviteCustomer_FullMethodName,
		customerv1.CustomerService_CreateCustomer_FullMethodName,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authinterceptor.JWTAuthInterceptor(publicKey, exemptMethods)),
	)

	customerv1.RegisterCustomerServiceServer(grpcServer, srv)
	grpc_health_v1.RegisterHealthServer(grpcServer, &grpcadapter.HealthChecker{})
	reflection.Register(grpcServer)

	addr := fmt.Sprintf(":%d", app.Config.GetGRPCPort())
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	app.CustomerServer = srv
	app.GrpcServer = grpcServer
	app.GrpcListener = lis
	return nil
}

// RunGRPCServer starts the gRPC server in an errgroup.
func (app *App) RunGRPCServer(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		app.Logger.Info("customer-service gRPC starting",
			zap.String("addr", app.GrpcListener.Addr().String()),
		)
		if err := app.GrpcServer.Serve(app.GrpcListener); err != nil {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		app.Logger.Info("shutting down gRPC server")
		app.GrpcServer.GracefulStop()
		return nil
	})
}

// Recoverer runs a function with panic recovery.
func (app *App) Recoverer(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			app.Logger.Error("panic recovered",
				zap.Any("panic", r),
				zap.String("stack", string(debug.Stack())),
			)
			go app.Recoverer(fn)
		}
	}()
	fn()
}
