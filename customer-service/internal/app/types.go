package app

import (
	customerrepo "github.com/beabys/wms/customer-service/internal/infrastructure/persistence/repository"
	grpcdapter "github.com/beabys/wms/customer-service/internal/infrastructure/adapters/grpc"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
)

// App is the application struct.
type App struct {
	Config             *Config
	Logger             logger.Logger
	DB                 database.Database
	GrpcServer         *grpcdapter.GRPCServer
	AuthClient         *grpcdapter.AuthClient
	CustomerRepo       *customerrepo.CustomerRepository
	CompanyRepo        *customerrepo.CompanyRepository
	CompanyRoleRepo    *customerrepo.CompanyRoleRepository
}
