package app

import (
	"go.uber.org/zap"

	"github.com/beabys/wms/notification-service/internal/app/ports"
	notificationgrpc "github.com/beabys/wms/notification-service/internal/infrastructure/adapters/grpc"
)

// App is the Application Struct.
type App struct {
	Config     ports.AppConfig
	Logger     *zap.Logger
	GrpcServer *notificationgrpc.GRPCServer
}
