package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/beabys/wms/inbound-service/internal/app/ports"
	inboundgrpc "github.com/beabys/wms/inbound-service/internal/infrastructure/adapters/grpc"
)

// App is the Application Struct
type App struct {
	Config     ports.AppConfig
	Logger     *zap.Logger
	DB         *pgxpool.Pool
	GrpcServer *inboundgrpc.GRPCServer
}
