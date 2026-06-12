package app

import (
	"fmt"
	"strings"

	config "github.com/beabys/ayotl"
)

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() (*Config, error) {
	c := config.NewWithParams(&config.Params{
		Defaults: config.ConfigMap{
			"server.host":            "0.0.0.0",
			"server.port":            8081,
			"auth_service.grpc_host": "localhost",
			"auth_service.grpc_port": 50001,
			"log.level":              "debug",
			"cors.allowed_origins":   "*",
		},
		EnvAlias: config.ConfigEnvAlias{
			"SERVER_HOST":            "server.host",
			"SERVER_PORT":            "server.port",
			"AUTH_SERVICE_GRPC_HOST": "auth_service.grpc_host",
			"AUTH_SERVICE_GRPC_PORT": "auth_service.grpc_port",
			"LOG_LEVEL":              "log.level",
			"CORS_ALLOWED_ORIGINS":   "cors.allowed_origins",
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
