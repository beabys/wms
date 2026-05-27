package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/beabys/wms/inbound-service/internal/app/ports"
	"github.com/beabys/wms/inbound-service/internal/application/inbound/usecase"
	inboundgrpc "github.com/beabys/wms/inbound-service/internal/infrastructure/adapters/grpc"
	kafkaadapter "github.com/beabys/wms/inbound-service/internal/infrastructure/adapters/kafka"
	prepository "github.com/beabys/wms/inbound-service/internal/infrastructure/persistence/repository"
)

// New creates a new App instance.
func New() *App {
	return &App{}
}

// SetConfigs sets the app configuration.
func (app *App) SetConfigs(cfg ports.AppConfig) error {
	app.Config = cfg
	return cfg.LoadConfigs()
}

// Setup initializes the logger.
func (app *App) Setup(configs ports.AppConfig) error {
	if err := app.SetConfigs(configs); err != nil {
		return err
	}

	appConfig := configs.GetConfigs()

	var logger *zap.Logger
	var err error
	if appConfig.Service.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	app.Logger = logger

	return nil
}

// GetLogger returns the app logger.
func (app *App) GetLogger() *zap.Logger {
	return app.Logger
}

// SetDB sets the database pool.
func (app *App) SetDB(db *pgxpool.Pool) {
	app.DB = db
}

// SetGRPCServer sets up the gRPC server with all dependencies.
func (a *App) SetGRPCServer() {
	cfg := a.Config.GetConfigs()

	// Repositories
	inboundRepo := prepository.NewInboundRepository(a.DB)

	// Kafka producer — use noop publisher in dev when Kafka is not configured
	var eventPub usecase.EventPublisher
	if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Brokers[0] != "" {
		eventPub = kafkaadapter.NewEventPublisher(cfg.Kafka.Brokers)
	} else {
		eventPub = kafkaadapter.NewNoopEventPublisher()
	}

	// Application service
	svc := usecase.NewInboundService(inboundRepo, eventPub)

	// gRPC server
	grpcServer := inboundgrpc.NewGRPCServer().
		SetConfig(cfg).
		SetLogger(a.Logger).
		SetService(svc)

	a.GrpcServer = grpcServer
}
