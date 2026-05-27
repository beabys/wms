package app

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/beabys/wms/pkg/logger"
	"github.com/beabys/wms/login-service/pkg/config"
)

// Start initializes and runs the login-service gRPC server.
func Start(ctx context.Context, wg *errgroup.Group) error {
	application := New()
	cfg := config.New()
	if err := application.Setup(cfg); err != nil {
		return err
	}

	appCfg := application.Config.GetConfigs()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	zapLog, err := logger.New(env)
	if err != nil {
		return err
	}
	application.SetLogger(zapLog)

	ctxDB := context.Background()
	if err := application.SetDBPool(ctxDB, appCfg.DBDsn); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	if err := application.SetGRPCServer(); err != nil {
		return fmt.Errorf("setup gRPC server: %w", err)
	}

	application.GrpcServer.Run(ctx, wg)
	return nil
}
