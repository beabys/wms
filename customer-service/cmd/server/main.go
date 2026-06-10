package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/beabys/wms/customer-service/internal/app"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stopFn := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stopFn()

	cfg, err := app.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	a := app.New()

	if err := a.Setup(cfg); err != nil {
		slog.Error("failed to setup application", "error", err)
		os.Exit(1)
	}

	wg, ctx := errgroup.WithContext(ctx)

	// Run gRPC server via the handler's Run method
	if a.GrpcServer != nil {
		a.GrpcServer.Run(ctx, wg)
	}

	if err := wg.Wait(); err != nil {
		a.Logger.Error("application stopped with error:", err)
	}

	// Graceful shutdown
	a.Shutdown()
	a.Logger.Info("application stopped")
}
