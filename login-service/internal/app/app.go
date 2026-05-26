package app

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"

	"github.com/beabys/wms/pkg/authinterceptor"
	pkgConfig "github.com/beabys/wms/login-service/pkg/config"

	"github.com/beabys/wms/login-service/internal/app/ports"
	"github.com/beabys/wms/login-service/internal/application/auth/repository"
	"github.com/beabys/wms/login-service/internal/application/auth/usecase"
	adaptergrpc "github.com/beabys/wms/login-service/internal/infrastructure/adapters/grpc"
	repoimpl "github.com/beabys/wms/login-service/internal/infrastructure/persistence/repository"
)

// New creates a new App instance.
func New() *App {
	return &App{}
}

// SetConfigs loads and sets the application configuration.
func (app *App) SetConfigs(ac ports.AppConfig) error {
	app.Config = ac
	return ac.LoadConfigs()
}

// SetLogger sets the application logger.
func (app *App) SetLogger(logger *zap.Logger) {
	app.Logger = logger
}

// Setup loads config, creates logger, and creates DB pool.
func (app *App) Setup(configs ports.AppConfig) error {
	if err := app.SetConfigs(configs); err != nil {
		return fmt.Errorf("set configs: %w", err)
	}

	return nil
}

// SetDBPool creates and sets the database connection pool.
func (app *App) SetDBPool(ctx context.Context, dsn string) error {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse pool config: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	app.Pool = pool
	return nil
}

// SetGRPCServer wires repository → usecase → gRPC server,
// registers auth service, user service, health, and reflection.
func (app *App) SetGRPCServer() error {
	cfg := app.Config.GetConfigs()

	// Parse keys
	privateKey, err := pkgConfig.ParsePrivateKey(cfg.JWTPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	publicKey, err := pkgConfig.ParsePublicKey(cfg.JWTPublicKeyPath)
	if err != nil {
		return fmt.Errorf("load public key: %w", err)
	}

	publicKeyPEM, err := pkgConfig.ReadPublicKeyPEM(cfg.JWTPublicKeyPath)
	if err != nil {
		return fmt.Errorf("read public key PEM: %w", err)
	}

	// Parse TTLs
	accessTTL, _ := time.ParseDuration(cfg.AccessTokenTTL)
	refreshTTL, _ := time.ParseDuration(cfg.RefreshTokenTTL)
	if accessTTL == 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL == 0 {
		refreshTTL = 7 * 24 * time.Hour
	}

	// Create repository
	var userRepo repository.UserRepository
	userRepo = repoimpl.NewUserRepo(app.Pool)

	// Create usecase
	authService := usecase.NewAuthService(userRepo, privateKey, publicKey, accessTTL, refreshTTL)

	// Exempt methods from JWT validation
	exemptMethods := []string{
		authv1.AuthService_Login_FullMethodName,
		authv1.AuthService_RefreshToken_FullMethodName,
		authv1.AuthService_ValidateToken_FullMethodName,
		authv1.AuthService_GetPublicKey_FullMethodName,
		authv1.UserService_CreateUser_FullMethodName,
	}

	// Create gRPC server with auth interceptor
	rpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authinterceptor.JWTAuthInterceptor(publicKey, exemptMethods)),
	)

	// Register services
	authv1.RegisterAuthServiceServer(rpcServer,
		adaptergrpc.NewAuthServer(authService, app.Logger, publicKeyPEM))
	authv1.RegisterUserServiceServer(rpcServer,
		adaptergrpc.NewUserServer(authService, app.Logger))

	// Register health check
	grpc_health_v1.RegisterHealthServer(rpcServer, &adaptergrpc.HealthChecker{})

	// Enable reflection for debugging
	reflection.Register(rpcServer)

	address := fmt.Sprintf(":%d", cfg.GRPCPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	app.GrpcServer = &adaptergrpc.GRPCServer{
		Server:   rpcServer,
		Listener: listener,
		Logger:   app.Logger,
	}

	app.Logger.Info("gRPC server setup complete", zap.Int("port", cfg.GRPCPort))
	return nil
}
