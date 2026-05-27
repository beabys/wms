package app

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	"github.com/beabys/wms/notification-service/internal/app/ports"
	"github.com/beabys/wms/notification-service/internal/application/notification/usecase"
	senderadapter "github.com/beabys/wms/notification-service/internal/infrastructure/adapters/sender"
	notificationgrpc "github.com/beabys/wms/notification-service/internal/infrastructure/adapters/grpc"
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

	cfgPath := "notification-service/config.yaml"
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

	if port := os.Getenv("GRPC_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.GRPC.Port)
	}

	return cfg, nil
}

// Start initializes and runs the notification-service gRPC server.
func Start(ctx context.Context, wg *errgroup.Group) error {
	a := New()

	config := &configLoader{}
	if err := a.Setup(config); err != nil {
		return err
	}

	cfg := config.GetConfigs()

	// Create sender based on config
	var sender usecase.Sender
	switch cfg.Sender.Type {
	case "console":
		sender = senderadapter.NewConsoleSender(a.Logger)
	case "smtp":
		sender = senderadapter.NewSMTPSender(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.Username,
			cfg.SMTP.Password,
		)
	default:
		return fmt.Errorf("unknown sender type: %s", cfg.Sender.Type)
	}

	// Application service
	svc := usecase.NewNotificationService(sender)

	// gRPC server
	grpcServer := notificationgrpc.NewGRPCServer().
		SetConfig(cfg).
		SetLogger(a.Logger).
		SetService(svc)

	a.GrpcServer = grpcServer

	a.GrpcServer.Run(ctx, wg)
	return nil
}
