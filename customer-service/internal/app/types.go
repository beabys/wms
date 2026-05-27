package app

import (
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/beabys/wms/customer-service/internal/app/ports"
	grpcadapter "github.com/beabys/wms/customer-service/internal/infrastructure/adapters/grpc"
)

// App is the Application Struct
type App struct {
	Config         ports.AppConfig
	Logger         *zap.Logger
	Pool           *pgxpool.Pool
	CustomerServer *grpcadapter.Server
	GrpcServer     *grpc.Server
	GrpcListener   net.Listener
}
