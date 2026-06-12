package app

import (
	usrrepo "github.com/beabys/wms/auth-service/internal/application/auth/usecase"
	grpcdapter "github.com/beabys/wms/auth-service/internal/infrastructure/adapters/grpc"
	authrepo "github.com/beabys/wms/auth-service/internal/infrastructure/persistence/repository"
	"github.com/beabys/wms/pkg/database"
	"github.com/beabys/wms/pkg/logger"
)

// App is the application struct.
type App struct {
	Config          *Config
	Logger          logger.Logger
	DB              database.Database
	GrpcServer      *grpcdapter.GRPCServer
	JWTService      *usrrepo.Service
	UserRepo        *authrepo.UserRepository
	AuthRepo        *authrepo.AuthRepository
	InviteRepo      *authrepo.InviteRepository
}
