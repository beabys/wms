package httpadapter

import (
	"context"
	"net/http"
	"time"

	"github.com/beabys/wms/pkg/logger"
	"golang.org/x/sync/errgroup"
)

// HTTPServer wraps an http.Server with lifecycle management.
type HTTPServer struct {
	Server *http.Server
	Logger logger.Logger
}

// NewHTTPServer creates a new HTTPServer wrapper.
func NewHTTPServer(addr string, handler http.Handler, log logger.Logger) *HTTPServer {
	return &HTTPServer{
		Server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Logger: log,
	}
}

// Run starts the HTTP server and handles graceful shutdown via the errgroup.
func (hs *HTTPServer) Run(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		hs.Logger.Info("HTTP server starting",
			logger.LogField{Key: "addr", Value: hs.Server.Addr},
		)
		if err := hs.Server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return nil
			}
			hs.Logger.Error("HTTP server stopped with error", err)
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		hs.Logger.Info("shutting down HTTP server gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := hs.Server.Shutdown(shutdownCtx); err != nil {
			hs.Logger.Error("error shutting server down", err)
			return err
		}
		return nil
	})
}
