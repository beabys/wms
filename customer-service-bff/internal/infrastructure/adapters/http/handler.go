package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/beabys/wms/customer-service-bff/internal/application"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// HttpServer is the HTTP server for the customer BFF.
type HttpServer struct {
	Server         *http.Server
	Config         *Config
	Logger         *zap.Logger
	CustomerClient application.CustomerService
}

// Config holds the HTTP server configuration.
type Config struct {
	Port int
}

// NewHttpServer creates a new HttpServer.
func NewHttpServer(cfg *Config, logger *zap.Logger, client application.CustomerService) *HttpServer {
	return &HttpServer{
		Config:         cfg,
		Logger:         logger,
		CustomerClient: client,
	}
}

// Run starts the HTTP server and handles graceful shutdown.
func (hs *HttpServer) Run(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		hs.Logger.Info("http server started", zap.Int("port", hs.Config.Port))
		if err := hs.Server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				return nil
			}
			hs.Logger.Error("http server stopped with error", zap.Error(err))
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		hs.Logger.Info("shutting down gracefully http server")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		if err := hs.Server.Shutdown(ctxTimeout); err != nil {
			hs.Logger.Error("error shutting server down", zap.Error(err))
			return err
		}
		return nil
	})
}

// SetServer sets the underlying http.Server. Used after muxhandler creates it.
func (hs *HttpServer) SetServer(srv *http.Server) {
	hs.Server = srv
}

// Addr returns the address the server listens on.
func (hs *HttpServer) Addr() string {
	if hs.Server != nil {
		return hs.Server.Addr
	}
	return fmt.Sprintf(":%d", hs.Config.Port)
}
