package grpcdapter

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func TestGRPCServerNew(t *testing.T) {
	gs := &GRPCServer{
		Server: grpc.NewServer(),
		Listener: func() net.Listener {
			l, _ := net.Listen("tcp", "127.0.0.1:0")
			return l
		}(),
	}
	assert.NotNil(t, gs.Server)
	assert.NotNil(t, gs.Listener)
}

func TestGRPCServerRunAndShutdown(t *testing.T) {
	log, err := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	gs := &GRPCServer{
		Server:   grpc.NewServer(),
		Listener: listener,
		Logger:   log,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	gs.Run(ctx, wg)

	// Wait briefly then cancel
	time.Sleep(50 * time.Millisecond)
	cancel()

	err = wg.Wait()
	assert.NoError(t, err)
}

func TestGRPCServerRunError(t *testing.T) {
	log, err := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	// Use a closed listener to trigger Serve error
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	listener.Close() // close immediately

	gs := &GRPCServer{
		Server:   grpc.NewServer(),
		Listener: listener,
		Logger:   log,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	gs.Run(ctx, wg)

	time.Sleep(50 * time.Millisecond)
	cancel()

	// Serve on closed listener should error, but our handler runs in goroutine
	// The error happens before shutdown, but the error group will capture it
	err = wg.Wait()
	// Expect an error from Serve on closed listener
	assert.Error(t, err)
}
