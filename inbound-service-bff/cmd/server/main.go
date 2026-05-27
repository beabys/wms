package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/beabys/wms/inbound-service-bff/pkg/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	wg, ctx := errgroup.WithContext(ctx)

	if err := app.Start(ctx, wg); err != nil {
		panic(err)
	}
	wg.Wait()
}
