package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the login-service-bff configuration.
type Config struct {
	GRPCAddress      string `yaml:"grpc_address"`
	HTTPPort         int    `yaml:"http_port"`
	AuthServiceAPIKey string `yaml:"auth_service_api_key"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		GRPCAddress: "localhost:50001",
		HTTPPort:    8080,
	}
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}
