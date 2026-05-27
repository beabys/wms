package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	internalapp "github.com/beabys/wms/notification-service/internal/app"
)

// Start initializes and runs the notification service.
func Start(ctx context.Context, wg *errgroup.Group) error {
	return internalapp.Start(ctx, wg)
}
