package app

import (
	"fmt"
	"strings"

	"github.com/beabys/ayotl"
)

// Config holds application configuration.
type Config struct {
	Server          ServerConfig          `mapstructure:"server"`
	CustomerService CustomerServiceConfig `mapstructure:"customer_service"`
	Log             LogConfig             `mapstructure:"log"`
	CORS            CORSConfig            `mapstructure:"cors"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// CustomerServiceConfig holds gRPC connection configuration for the customer service.
type CustomerServiceConfig struct {
	GRPCHost string `mapstructure:"grpc_host"`
	GRPCPort int    `mapstructure:"grpc_port"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// CORSConfig holds CORS configuration.
type CORSConfig struct {
	AllowedOrigins string `mapstructure:"allowed_origins"`
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() (*Config, error) {
	c := config.NewWithParams(&config.Params{
		Defaults: config.ConfigMap{
			"server.host":                  "0.0.0.0",
			"server.port":                  8082,
			"customer_service.grpc_host":   "localhost",
			"customer_service.grpc_port":   50002,
			"log.level":                    "debug",
			"cors.allowed_origins":         "*",
		},
		EnvAlias: config.ConfigEnvAlias{
			"SERVER_HOST":                  "server.host",
			"SERVER_PORT":                  "server.port",
			"CUSTOMER_SERVICE_GRPC_HOST":   "customer_service.grpc_host",
			"CUSTOMER_SERVICE_GRPC_PORT":   "customer_service.grpc_port",
			"LOG_LEVEL":                    "log.level",
			"CORS_ALLOWED_ORIGINS":         "cors.allowed_origins",
		},
	})

	c.WithEnv()
	c.LoadConfigs()

	var cfg Config
	if err := c.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Origins returns the allowed origins, defaulting to ["*"] if empty.
func (c *CORSConfig) Origins() []string {
	val := c.AllowedOrigins
	if val == "" {
		return []string{"*"}
	}
	parts := strings.Split(val, ",")
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = strings.TrimSpace(p)
	}
	return result
}

// Address returns the HTTP server address.
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LevelStr returns the slog-compatible log level.
func (c *LogConfig) LevelStr() string {
	if c.Level == "" {
		return "debug"
	}
	return c.Level
}
