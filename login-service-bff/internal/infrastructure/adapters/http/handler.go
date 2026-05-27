package http

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/beabys/wms/login-service-bff/internal/app/config"

	authv1 "github.com/beabys/wms/proto/gen/go/auth/v1"
)

// GRPCAuthClient defines the interface the HTTP handler needs from the gRPC layer.
// Defined by the consumer (http) per Go interface best practices.
type GRPCAuthClient interface {
	Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error)
	Register(ctx context.Context, email, password, companyName string) (*authv1.CreateUserResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*authv1.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error)
	GetUser(ctx context.Context, userID string, token string) (*authv1.GetUserResponse, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	Close() error
}

// HttpServer holds all dependencies for the BFF HTTP server.
type HttpServer struct {
	Server     *http.Server
	Config     *config.Config
	Logger     *zap.Logger
	GRPCClient GRPCAuthClient
}

// NewHttpServer creates a new HttpServer.
func NewHttpServer() *HttpServer {
	return &HttpServer{}
}

// SetConfig sets the config and returns the server for chaining.
func (hs *HttpServer) SetConfig(c *config.Config) *HttpServer {
	hs.Config = c
	return hs
}

// SetLogger sets the logger and returns the server for chaining.
func (hs *HttpServer) SetLogger(l *zap.Logger) *HttpServer {
	hs.Logger = l
	return hs
}

// SetGRPCClient sets the gRPC client and returns the server for chaining.
func (hs *HttpServer) SetGRPCClient(c GRPCAuthClient) *HttpServer {
	hs.GRPCClient = c
	return hs
}

// Run starts the HTTP server and handles graceful shutdown.
func (hs *HttpServer) Run(ctx context.Context, wg *errgroup.Group) {
	wg.Go(func() error {
		hs.Logger.Info("http server started", zap.Int("port", hs.Config.HTTPPort))
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
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		if err := hs.Server.Shutdown(shutdownCtx); err != nil {
			hs.Logger.Error("error shutting server down", zap.Error(err))
			return err
		}
		return nil
	})
}
