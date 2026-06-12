package httpadapter

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

func TestNewHTTPServer(t *testing.T) {
	log, err := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := NewHTTPServer("127.0.0.1:0", handler, log)
	require.NotNil(t, srv)
	assert.NotNil(t, srv.Server)
	assert.Equal(t, "127.0.0.1:0", srv.Server.Addr)
	assert.Equal(t, 15*time.Second, srv.Server.ReadTimeout)
	assert.Equal(t, 15*time.Second, srv.Server.WriteTimeout)
	assert.Equal(t, 60*time.Second, srv.Server.IdleTimeout)
}

func TestHTTPServerRunAndShutdown(t *testing.T) {
	log, err := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := NewHTTPServer("127.0.0.1:0", handler, log)

	// Change to random port so we don't conflict
	srv.Server.Addr = "127.0.0.1:0"

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	srv.Run(ctx, wg)

	// Give server time to start
	time.Sleep(50 * time.Millisecond)
	cancel()

	err = wg.Wait()
	assert.NoError(t, err)
}

func TestHTTPServerDefaults(t *testing.T) {
	log, err := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	require.NoError(t, err)

	srv := NewHTTPServer(":3000", http.NewServeMux(), log)
	assert.Equal(t, ":3000", srv.Server.Addr)
	assert.Equal(t, 15*time.Second, srv.Server.ReadTimeout)
	assert.Equal(t, 15*time.Second, srv.Server.WriteTimeout)
	assert.Equal(t, 60*time.Second, srv.Server.IdleTimeout)
}
