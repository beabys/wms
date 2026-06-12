package app

import (
	"github.com/beabys/wms/customer-service-bff/internal/application/customer/usecase"
	grpcdapter "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/grpc"
	httpadapter "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http"
	"github.com/beabys/wms/pkg/logger"
)

// App is the application struct holding all dependencies.
type App struct {
	Config     *Config
	Logger     logger.Logger
	HTTPServer *httpadapter.HTTPServer
	Client     *grpcdapter.CustomerClient
	UseCase    *usecase.CustomerUseCase
}
