// Package app provides the public Start function for the inbound-service-bff.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	grpcadapter "github.com/beabys/wms/inbound-service-bff/internal/infrastructure/adapters/grpc"
	httpadapter "github.com/beabys/wms/inbound-service-bff/internal/infrastructure/adapters/http"
	"github.com/beabys/wms/pkg/logger"
)

// Config holds BFF application configuration.
type Config struct {
	Service struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
	} `yaml:"service"`
	HTTP struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	InboundService struct {
		GRPCAddr string `yaml:"grpc_addr"`
	} `yaml:"inbound_service"`
}

func loadConfig() *Config {
	cfg := &Config{}

	cfgPath := "inbound-service-bff/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		cfgPath = envPath
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		logger.MustNew("development").Fatal("read config", zap.Error(err))
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		logger.MustNew("development").Fatal("parse config", zap.Error(err))
	}

	// Override from env
	if port := os.Getenv("HTTP_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.HTTP.Port)
	}
	if addr := os.Getenv("INBOUND_GRPC_ADDR"); addr != "" {
		cfg.InboundService.GRPCAddr = addr
	}

	return cfg
}

// Start initializes and runs the inbound-service-bff HTTP server.
func Start(ctx context.Context, wg *errgroup.Group) error {
	cfg := loadConfig()

	log, err := logger.New(cfg.Service.Env)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	grpcClient, err := grpcadapter.NewClient(cfg.InboundService.GRPCAddr)
	if err != nil {
		return fmt.Errorf("create grpc client: %w", err)
	}

	httpServer := httpadapter.NewHttpServer().
		SetConfig(&httpadapter.Config{Port: cfg.HTTP.Port}).
		SetLogger(log).
		SetInboundClient(grpcClient)

	mux, err := httpadapter.NewMuxHandler(httpServer)
	if err != nil {
		return fmt.Errorf("create mux handler: %w", err)
	}

	httpServer.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: mux,
	}

	httpServer.Run(ctx, wg)
	return nil
}
