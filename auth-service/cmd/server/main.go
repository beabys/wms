package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/beabys/wms/auth-service/internal/app"
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

	if a.GrpcServer != nil {
		a.GrpcServer.Run(ctx, wg)
	}

	if err := wg.Wait(); err != nil {
		a.Logger.Error("application stopped with error:", err)
	}
	a.Logger.Info("application stopped")
}
