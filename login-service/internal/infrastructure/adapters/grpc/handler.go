package grpc

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// GRPCServer is a wrapper around *grpc.Server with lifecycle management.
type GRPCServer struct {
	Server   *grpc.Server
	Listener net.Listener
	Logger   *zap.Logger
}

// Run starts the gRPC server and handles graceful shutdown.
func (gs *GRPCServer) Run(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		gs.Logger.Info("gRPC server starting")
		if err := gs.Server.Serve(gs.Listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			gs.Logger.Error("gRPC server failed", zap.Error(err))
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		gs.Logger.Info("shutting down gRPC server")
		gs.Server.GracefulStop()
		return nil
	})
}
