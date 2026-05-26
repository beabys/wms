// Package app provides the public Start function for the customer-service-bff.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	bffgrpc "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/grpc"
	bffhttp "github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http"
	"github.com/beabys/wms/pkg/logger"
)

// Config holds BFF application configuration.
type Config struct {
	CustomerGRPCAddress string `yaml:"customer_grpc_address"`
	HTTPPort            int    `yaml:"http_port"`
	LogLevel            string `yaml:"log_level"`
	JWTPublicKeyPath    string `yaml:"jwt_public_key_path"`
	ServiceAPIKey       string `yaml:"service_api_key"`
}

func loadConfig() Config {
	cfg := Config{}

	data, err := os.ReadFile("customer-service-bff/config.yaml")
	if err != nil {
		logger.MustNew("development").Fatal("read config", zap.Error(err))
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.MustNew("development").Fatal("parse config", zap.Error(err))
	}

	// Override from env
	if port := os.Getenv("HTTP_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.HTTPPort)
	}
	if addr := os.Getenv("CUSTOMER_GRPC_ADDR"); addr != "" {
		cfg.CustomerGRPCAddress = addr
	}

	return cfg
}

// Start initializes and runs the customer-service-bff HTTP server.
func Start(ctx context.Context, wg *errgroup.Group) error {
	cfg := loadConfig()

	log, err := logger.New("development")
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	conn, err := grpc.DialContext(ctx, cfg.CustomerGRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("connect to customer service: %w", err)
	}

	grpcClient := bffgrpc.NewClient(conn)

	httpCfg := &bffhttp.Config{
		Port: cfg.HTTPPort,
	}
	httpServer := bffhttp.NewHttpServer(httpCfg, log, grpcClient)

	handler, err := bffhttp.NewMuxHandler(httpServer)
	if err != nil {
		return fmt.Errorf("create mux handler: %w", err)
	}

	httpServer.SetServer(&http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: handler,
	})

	httpServer.Run(ctx, wg)
	return nil
}
