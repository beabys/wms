package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	customerapp "github.com/beabys/wms/customer-service/pkg/app"
	inboundapp "github.com/beabys/wms/inbound-service/pkg/app"
	inventoryapp "github.com/beabys/wms/inventory-service/pkg/app"
	loginapp "github.com/beabys/wms/login-service/pkg/app"
	notificationapp "github.com/beabys/wms/notification-service/pkg/app"

	customerbffapp "github.com/beabys/wms/customer-service-bff/pkg/app"
	inboundbffapp "github.com/beabys/wms/inbound-service-bff/pkg/app"
	inventorybffapp "github.com/beabys/wms/inventory-service-bff/pkg/app"
	loginbffapp "github.com/beabys/wms/login-service-bff/pkg/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	wg, ctx := errgroup.WithContext(ctx)

	// Domain services (gRPC)
	if err := loginapp.Start(ctx, wg); err != nil {
		panic(err)
	}
	if err := customerapp.Start(ctx, wg); err != nil {
		panic(err)
	}
	if err := inboundapp.Start(ctx, wg); err != nil {
		panic(err)
	}
	if err := inventoryapp.Start(ctx, wg); err != nil {
		panic(err)
	}
	if err := notificationapp.Start(ctx, wg); err != nil {
		panic(err)
	}

	// BFFs (HTTP) — connect to domain services via localhost
	if err := loginbffapp.Start(ctx, wg); err != nil {
		panic(err)
	}
	if err := customerbffapp.Start(ctx, wg); err != nil {
		panic(err)
	}
	if err := inboundbffapp.Start(ctx, wg); err != nil {
		panic(err)
	}
	if err := inventorybffapp.Start(ctx, wg); err != nil {
		panic(err)
	}

	wg.Wait()
}
