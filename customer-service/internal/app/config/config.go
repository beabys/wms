package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the customer-service configuration.
type Config struct {
	DBDSN            string `yaml:"db_dsn"`
	GRPCPort         int    `yaml:"grpc_port"`
	LogLevel         string `yaml:"log_level"`
	JWTPublicKeyPath string `yaml:"jwt_public_key_path"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		GRPCPort: 50002,
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
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadConfig implements ports.AppConfig.
func (c *Config) LoadConfig(path string) error {
	cfg, err := LoadConfig(path)
	if err != nil {
		return err
	}
	*c = *cfg
	return nil
}

// GetGRPCPort implements ports.AppConfig.
func (c *Config) GetGRPCPort() int { return c.GRPCPort }

// GetDBDSN implements ports.AppConfig.
func (c *Config) GetDBDSN() string { return c.DBDSN }

// GetLogLevel implements ports.AppConfig.
func (c *Config) GetLogLevel() string { return c.LogLevel }

// GetJWTPublicKeyPath implements ports.AppConfig.
func (c *Config) GetJWTPublicKeyPath() string { return c.JWTPublicKeyPath }
