package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/beabys/wms/customer-service-bff/internal/app"
	"github.com/beabys/wms/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stopFn := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopFn()

	cfg, err := app.LoadConfig()
	if err != nil {
		panic(err)
	}

	a := app.New()
	if err := a.Setup(cfg); err != nil {
		panic(err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	a.HTTPServer.Run(ctx, wg)

	a.Logger.Info("service started",
		logger.LogField{Key: "http_addr", Value: cfg.Server.Host},
		logger.LogField{Key: "http_port", Value: cfg.Server.Port},
		logger.LogField{Key: "grpc_target", Value: cfg.CustomerService.GRPCHost},
	)

	if err := wg.Wait(); err != nil {
		a.Logger.Error("service stopped with error", err)
		os.Exit(1)
	}
	a.Logger.Info("service stopped gracefully")
}
