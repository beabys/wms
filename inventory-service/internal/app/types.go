package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/beabys/wms/inventory-service/internal/app/ports"
	inventorygrpc "github.com/beabys/wms/inventory-service/internal/infrastructure/adapters/grpc"
	kafkaadapter "github.com/beabys/wms/inventory-service/internal/infrastructure/adapters/kafka"
)

// App is the Application Struct
type App struct {
	Config        ports.AppConfig
	Logger        *zap.Logger
	DB            *pgxpool.Pool
	GrpcServer    *inventorygrpc.GRPCServer
	KafkaConsumer *kafkaadapter.InboundEventConsumer
}
