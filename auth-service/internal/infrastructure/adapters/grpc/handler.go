package grpcdapter

import (
	"context"
	"errors"
	"net"

	"github.com/beabys/wms/pkg/logger"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// GRPCServer wraps the gRPC server with lifecycle management.
type GRPCServer struct {
	Server   *grpc.Server
	Listener net.Listener
	Logger   logger.Logger
}

// Run starts the gRPC server and handles graceful shutdown.
func (gs *GRPCServer) Run(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		gs.Logger.Info("gRPC server started", logger.LogField{Key: "address", Value: gs.Listener.Addr().String()})
		if err := gs.Server.Serve(gs.Listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			gs.Logger.Error("gRPC server failed to serve", err)
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		gs.Logger.Info("shutting down gRPC server gracefully")
		gs.Server.GracefulStop()
		return nil
	})
}
