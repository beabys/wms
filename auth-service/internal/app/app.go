package app

import (
	"fmt"
	"net"

	"github.com/beabys/wms/auth-service/internal/application/auth/repository"
	"github.com/beabys/wms/auth-service/internal/application/auth/usecase"
	grpcdapter "github.com/beabys/wms/auth-service/internal/infrastructure/adapters/grpc"
	authrepo "github.com/beabys/wms/auth-service/internal/infrastructure/persistence/repository"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
)

// New creates a new App.
func New() *App {
	return &App{}
}

// Setup initialises the application dependencies.
func (a *App) Setup(cfg *Config) error {
	// 1. Store config
	a.Config = cfg

	// 2. Init logger
	logLevel := zapcore.DebugLevel
	switch cfg.Log.LevelStr() {
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	}
	log, err := logger.NewZapLogger([]string{}, []string{}, logLevel)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	a.Logger = log
	a.Logger.Info("config loaded", logger.LogField{Key: "server_addr", Value: cfg.Server.Address()})

	// 3. Connect to PostgreSQL
	db := database.New()
	dbConfig := &database.PostgresConfig{
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetimeDuration(),
		ConnectionRetries: 3,
	}
	db.SetConfigs(dbConfig)
	if err := db.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.DB = db
	a.Logger.Info("database connected")

	// 4. Init JWT service
	jwtSvc, err := usecase.NewService(
		cfg.JWT.PrivateKeyPEM,
		cfg.JWT.PublicKeyPEM,
		cfg.JWT.AccessTokenTTLDuration(),
		cfg.JWT.RefreshTokenTTLDuration(),
	)
	if err != nil {
		return fmt.Errorf("failed to init JWT service: %w", err)
	}
	a.JWTService = jwtSvc

	// 5. Init repositories
	userRepo := authrepo.NewUserRepository(a.Logger, a.DB)
	authRepo := authrepo.NewAuthRepository(a.Logger, a.DB)
	inviteRepo := authrepo.NewInviteRepository(a.Logger, a.DB)
	a.UserRepo = userRepo
	a.AuthRepo = authRepo
	a.InviteRepo = inviteRepo

	// 6. Init use cases
	var userRepoIface repository.UserRepository = userRepo
	var authRepoIface repository.AuthRepository = authRepo
	var inviteRepoIface repository.InviteRepository = inviteRepo

	loginUC := usecase.NewLoginUseCase(a.Logger, userRepoIface, authRepoIface, jwtSvc)
	userUC := usecase.NewUserUseCase(a.Logger, userRepoIface)
	inviteUC := usecase.NewInviteUseCase(a.Logger, inviteRepoIface, userRepoIface)

	// 7. Set up gRPC server with interceptor
	exemptMethods := grpcdapter.GetPublicKeyExemptMethods()
	authInterceptor := grpcdapter.NewAuthInterceptor(jwtSvc, exemptMethods, a.Logger)

	recoveryOpt := grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		a.Logger.Error("panic recovered", fmt.Errorf("%v", p))
		return status.Errorf(codes.Internal, "internal server error")
	})

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(recoveryOpt),
			grpcdapter.LoggingInterceptor(a.Logger),
			authInterceptor.Unary(),
		),
	)

	authServer := grpcdapter.NewAuthServer(loginUC, inviteUC)
	userServer := grpcdapter.NewUserServer(userUC)
	healthServer := grpcdapter.NewHealthServer()

	authv1.RegisterAuthServiceServer(grpcServer, authServer)
	authv1.RegisterUserServiceServer(grpcServer, userServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.Server.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", cfg.Server.Address(), err)
	}

	gs := &grpcdapter.GRPCServer{
		Server:   grpcServer,
		Listener: listener,
		Logger:   a.Logger,
	}
	a.GrpcServer = gs

	a.Logger.Info("gRPC server configured with auth interceptor",
		logger.LogField{Key: "exempt_methods", Value: exemptMethods},
	)

	return nil
}

// Run starts the gRPC server.
func (a *App) Run() error {
	if a.GrpcServer == nil {
		return fmt.Errorf("grpc server not configured")
	}
	a.Logger.Info("starting gRPC server", logger.LogField{Key: "address", Value: a.GrpcServer.Listener.Addr().String()})
	return a.GrpcServer.Server.Serve(a.GrpcServer.Listener)
}

// Shutdown gracefully stops the gRPC server.
func (a *App) Shutdown() {
	if a.GrpcServer != nil {
		a.GrpcServer.Server.GracefulStop()
	}
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			a.Logger.Error("failed to close database", err)
		}
	}
}
