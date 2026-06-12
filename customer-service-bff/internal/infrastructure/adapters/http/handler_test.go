package httpadapter

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/beabys/wms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

func TestNewHTTPServer(t *testing.T) {
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	hs := NewHTTPServer("127.0.0.1:0", handler, log)
	assert.NotNil(t, hs)
	assert.NotNil(t, hs.Server)
}

func TestHTTPServerRunAndShutdown(t *testing.T) {
	log, _ := logger.NewZapLogger([]string{}, []string{}, zapcore.DebugLevel)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	hs := NewHTTPServer("127.0.0.1:0", handler, log)

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	hs.Run(ctx, wg)

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Trigger shutdown
	cancel()

	err := wg.Wait()
	assert.NoError(t, err)
}
