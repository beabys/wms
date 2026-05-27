package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/beabys/wms/inventory-service/internal/app/ports"
	"github.com/beabys/wms/inventory-service/internal/application/inventory/usecase"
	inventorygrpc "github.com/beabys/wms/inventory-service/internal/infrastructure/adapters/grpc"
	kafkaadapter "github.com/beabys/wms/inventory-service/internal/infrastructure/adapters/kafka"
	prepository "github.com/beabys/wms/inventory-service/internal/infrastructure/persistence/repository"
)

// ensure InventoryService is the concrete type we assert for Kafka consumer
var _ usecase.InventoryServiceHandler = (*usecase.InventoryService)(nil)

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
	productRepo := prepository.NewProductRepository(a.DB)
	stockRepo := prepository.NewStockRepository(a.DB)
	binLocationRepo := prepository.NewBinLocationRepository(a.DB)

	// Kafka producer — use noop publisher in dev when Kafka is not configured
	var eventPub usecase.EventPublisher
	if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Brokers[0] != "" {
		eventPub = kafkaadapter.NewEventPublisher(cfg.Kafka.Brokers)
	} else {
		eventPub = kafkaadapter.NewNoopEventPublisher()
	}

	// Application service
	svc := usecase.NewInventoryService(productRepo, stockRepo, binLocationRepo, eventPub)

	// gRPC server
	grpcServer := inventorygrpc.NewGRPCServer().
		SetConfig(cfg).
		SetLogger(a.Logger).
		SetService(svc)

	a.GrpcServer = grpcServer
}

// SetKafkaConsumer sets up the Kafka consumer if brokers are configured.
func (a *App) SetKafkaConsumer() {
	cfg := a.Config.GetConfigs()

	if len(cfg.Kafka.Brokers) == 0 || cfg.Kafka.Brokers[0] == "" {
		a.Logger.Info("Kafka consumer disabled — no brokers configured")
		return
	}

	// The service is stored as InventoryServiceHandler interface;
	// assert the concrete type needed by the Kafka consumer.
	svc, ok := a.GrpcServer.Service.(*usecase.InventoryService)
	if !ok {
		a.Logger.Warn("Kafka consumer skipped — service is not the expected concrete type")
		return
	}

	a.KafkaConsumer = kafkaadapter.NewInboundEventConsumer(cfg.Kafka.Brokers, svc, a.Logger)
	a.Logger.Info("Kafka consumer initialized", zap.Strings("brokers", cfg.Kafka.Brokers))
}
