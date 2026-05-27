package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/beabys/wms/login-service/internal/app/ports"
	adaptergrpc "github.com/beabys/wms/login-service/internal/infrastructure/adapters/grpc"
)

// App is the Application Struct.
type App struct {
	Config     ports.AppConfig
	Logger     *zap.Logger
	Pool       *pgxpool.Pool
	GrpcServer *adaptergrpc.GRPCServer
}
