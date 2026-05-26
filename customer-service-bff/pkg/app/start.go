package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	internalapp "github.com/beabys/wms/customer-service-bff/internal/app"
)

func Start(ctx context.Context, wg *errgroup.Group) error {
	return internalapp.Start(ctx, wg)
}
