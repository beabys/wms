package app

import (
	"fmt"

	"github.com/beabys/wms/customer-service-bff/internal/application/customer/usecase"
	grpcdapter "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/grpc"
	httpadapter "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http"
	"github.com/beabys/wms/pkg/logger"
	"go.uber.org/zap/zapcore"
)

// New creates a new App.
func New() *App {
	return &App{}
}

// Setup initialises all application dependencies.
func (a *App) Setup(cfg *Config) error {
	// 1. Init logger
	logLevel := zapcore.DebugLevel
	if cfg.Log.Level != "" {
		switch cfg.Log.Level {
		case "debug":
			logLevel = zapcore.DebugLevel
		case "info":
			logLevel = zapcore.InfoLevel
		case "warn":
			logLevel = zapcore.WarnLevel
		case "error":
			logLevel = zapcore.ErrorLevel
		}
	}
	log, err := logger.NewZapLogger([]string{}, []string{}, logLevel)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	a.Logger = log

	// 2. Dial gRPC client to customer-service
	grpcTarget := fmt.Sprintf("%s:%d", cfg.CustomerService.GRPCHost, cfg.CustomerService.GRPCPort)
	grpcClient, err := grpcdapter.NewCustomerClient(grpcTarget)
	if err != nil {
		return fmt.Errorf("grpc client: %w", err)
	}
	a.Client = grpcClient

	// 3. Create gRPC adapter (thin layer: proto in/out, error mapping)
	grpcAdapter := grpcdapter.NewCustomerServiceAdapter(grpcClient)

	// 4. Create use case (business logic)
	customerUseCase := usecase.New(log, grpcAdapter)
	a.UseCase = customerUseCase

	// 5. Create HTTP handlers
	server := httpadapter.NewServer(customerUseCase, log)

	// 6. Setup router with middleware
	allowedOrigins := cfg.CORS.Origins()
	handler, err := httpadapter.NewMuxHandler(server, log, allowedOrigins)
	if err != nil {
		return fmt.Errorf("mux handler: %w", err)
	}

	// 7. Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	a.HTTPServer = httpadapter.NewHTTPServer(addr, handler, log)

	return nil
}
