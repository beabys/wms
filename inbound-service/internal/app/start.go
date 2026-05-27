package app

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	"github.com/beabys/wms/inbound-service/internal/app/ports"
	"github.com/beabys/wms/pkg/postgres"
)

// configLoader loads config from config.yaml with env overrides.
type configLoader struct {
	cfg *ports.Config
}

func (c *configLoader) LoadConfigs() error {
	cfg, err := loadConfigFromYAML()
	if err != nil {
		return err
	}
	c.cfg = cfg
	return nil
}

func (c *configLoader) GetConfigs() *ports.Config {
	return c.cfg
}

func loadConfigFromYAML() (*ports.Config, error) {
	cfg := &ports.Config{}

	cfgPath := "inbound-service/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		cfgPath = envPath
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Override from env
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}
	if port := os.Getenv("GRPC_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.GRPC.Port)
	}

	return cfg, nil
}

// Start initializes and runs the inbound-service gRPC server.
func Start(ctx context.Context, wg *errgroup.Group) error {
	a := New()

	config := &configLoader{}
	if err := a.Setup(config); err != nil {
		return err
	}

	cfg := config.GetConfigs()
	pool, err := postgres.NewPool(ctx, cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	a.SetDB(pool)
	a.SetGRPCServer()

	a.GrpcServer.Run(ctx, wg)
	return nil
}
