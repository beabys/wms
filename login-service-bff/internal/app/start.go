// Package app provides the public Start function for the login-service-bff.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/beabys/wms/pkg/logger"

	appConfig "github.com/beabys/wms/login-service-bff/internal/app/config"
	bffgrpc "github.com/beabys/wms/login-service-bff/internal/infrastructure/adapters/grpc"
	bffhttp "github.com/beabys/wms/login-service-bff/internal/infrastructure/adapters/http"
)

// Start initializes and runs the login-service-bff HTTP server.
func Start(ctx context.Context, wg *errgroup.Group) error {
	cfg, err := appConfig.LoadConfig("config.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	zapLog, err := logger.New(env)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	grpcClient, err := bffgrpc.NewClient(ctx, cfg.GRPCAddress)
	if err != nil {
		return fmt.Errorf("create gRPC client: %w", err)
	}

	httpServer := bffhttp.NewHttpServer().
		SetConfig(cfg).
		SetLogger(zapLog).
		SetGRPCClient(grpcClient)

	muxHandler, err := bffhttp.NewMuxHandler(httpServer)
	if err != nil {
		return fmt.Errorf("create mux handler: %w", err)
	}

	httpServer.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: muxHandler,
	}

	httpServer.Run(ctx, wg)
	return nil
}
