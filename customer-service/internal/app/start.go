package app

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/beabys/wms/customer-service/internal/app/config"
)

// Start initializes and runs the customer-service gRPC server.
func Start(ctx context.Context, wg *errgroup.Group) error {
	application := New()
	cfg := config.DefaultConfig()
	if err := application.Setup(cfg); err != nil {
		return err
	}

	if err := application.ConnectDB(ctx); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	if err := application.SetGRPCServer(); err != nil {
		return fmt.Errorf("setup gRPC server: %w", err)
	}

	application.RunGRPCServer(ctx, wg)
	return nil
}
