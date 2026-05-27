package http

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	inboundv1 "github.com/beabys/wms/proto/gen/go/inbound/v1"
)

// GRPCInboundClient defines the interface the HTTP handler needs from the gRPC layer.
type GRPCInboundClient interface {
	CreateInbound(ctx context.Context, customerID, expectedDate, notes string, items []*inboundv1.InboundItem, token string) (*inboundv1.Inbound, error)
	GetInbound(ctx context.Context, id string, token string) (*inboundv1.Inbound, error)
	ListInbounds(ctx context.Context, customerID, status string, pageSize int32, pageToken string, token string) ([]*inboundv1.Inbound, error)
	InspectInbound(ctx context.Context, id, inspectorID, notes string, passed bool, token string) (*inboundv1.Inbound, error)
	ApproveInbound(ctx context.Context, id string, token string) (*inboundv1.Inbound, error)
	FlagInbound(ctx context.Context, id, reason string, token string) (*inboundv1.Inbound, error)
	HoldInbound(ctx context.Context, id, reason string, token string) (*inboundv1.Inbound, error)
	ReleaseInbound(ctx context.Context, id string, token string) (*inboundv1.Inbound, error)
	Close() error
}

// HttpServer holds dependencies for the HTTP server.
type HttpServer struct {
	Server        *http.Server
	Config        *Config
	Logger        *zap.Logger
	InboundClient GRPCInboundClient
}

// Config holds BFF HTTP configuration.
type Config struct {
	Port int `yaml:"port"`
}

// NewHttpServer creates a new HttpServer.
func NewHttpServer() *HttpServer {
	return &HttpServer{}
}

// SetConfig sets the server config.
func (hs *HttpServer) SetConfig(cfg *Config) *HttpServer {
	hs.Config = cfg
	return hs
}

// SetLogger sets the logger.
func (hs *HttpServer) SetLogger(l *zap.Logger) *HttpServer {
	hs.Logger = l
	return hs
}

// SetInboundClient sets the gRPC client for inbound operations.
func (hs *HttpServer) SetInboundClient(client GRPCInboundClient) *HttpServer {
	hs.InboundClient = client
	return hs
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

