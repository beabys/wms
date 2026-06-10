package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/beabys/ayotl"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	AuthService AuthServiceConfig `mapstructure:"auth_service"`
	Log         LogConfig         `mapstructure:"log"`
}

// ServerConfig holds gRPC server configuration.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig holds PostgreSQL connection configuration.
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

// JWTConfig holds JWT public key configuration for token validation.
type JWTConfig struct {
	PublicKeyPath string `mapstructure:"public_key_path"`
	PublicKeyPEM  string `mapstructure:"public_key"`
}

// AuthServiceConfig holds the auth-service gRPC connection configuration.
type AuthServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() (*Config, error) {
	c := config.NewWithParams(&config.Params{
		Defaults: config.ConfigMap{
			"server.host":                 "0.0.0.0",
			"server.port":                 50002,
			"database.host":               "localhost",
			"database.port":               5432,
			"database.user":               "postgres",
			"database.password":           "postgres",
			"database.dbname":             "wms_customer",
			"database.sslmode":            "disable",
			"database.max_idle_conns":     10,
			"database.max_open_conns":     25,
			"database.conn_max_lifetime":  "5m",
			"auth_service.host":           "localhost",
			"auth_service.port":           50001,
			"log.level":                   "debug",
		},
		EnvAlias: config.ConfigEnvAlias{
			"SERVER_HOST":               "server.host",
			"SERVER_PORT":               "server.port",
			"DB_HOST":                   "database.host",
			"DB_PORT":                   "database.port",
			"DB_USER":                   "database.user",
			"DB_PASSWORD":               "database.password",
			"DB_NAME":                   "database.dbname",
			"DB_SSLMODE":                "database.sslmode",
			"DB_MAX_IDLE_CONNS":         "database.max_idle_conns",
			"DB_MAX_OPEN_CONNS":         "database.max_open_conns",
			"DB_CONN_MAX_LIFETIME":      "database.conn_max_lifetime",
			"JWT_PUBLIC_KEY_PATH":       "jwt.public_key_path",
			"JWT_PUBLIC_KEY":            "jwt.public_key",
			"AUTH_SERVICE_HOST":         "auth_service.host",
			"AUTH_SERVICE_PORT":         "auth_service.port",
			"LOG_LEVEL":                 "log.level",
		},
	})

	c.WithEnv()
	c.LoadConfigs()

	var cfg Config
	if err := c.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Handle JWT public key from file
	if cfg.JWT.PublicKeyPath != "" && cfg.JWT.PublicKeyPEM == "" {
		data, err := os.ReadFile(cfg.JWT.PublicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read public key: %w", err)
		}
		cfg.JWT.PublicKeyPEM = string(data)
	}

	if cfg.JWT.PublicKeyPEM == "" {
		return nil, fmt.Errorf("JWT public key is required: set JWT_PUBLIC_KEY_PATH or JWT_PUBLIC_KEY env vars")
	}

	return &cfg, nil
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// ConnMaxLifetimeDuration parses ConnMaxLifetime as a duration.
func (c *DatabaseConfig) ConnMaxLifetimeDuration() time.Duration {
	d, err := time.ParseDuration(c.ConnMaxLifetime)
	if err != nil {
		return 5 * time.Minute
	}
	return d
}

// Address returns the gRPC server address.
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Address returns the auth-service gRPC address.
func (c *AuthServiceConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LevelStr returns the log level string.
func (c *LogConfig) LevelStr() string {
	return strings.ToLower(c.Level)
}
