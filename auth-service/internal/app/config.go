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
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
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

// JWTConfig holds JWT and RSA key configuration.
type JWTConfig struct {
	PrivateKeyPath  string `mapstructure:"private_key_path"`
	PublicKeyPath   string `mapstructure:"public_key_path"`
	PrivateKeyPEM   string `mapstructure:"private_key"`
	PublicKeyPEM    string `mapstructure:"public_key"`
	AccessTokenTTL  string `mapstructure:"access_token_ttl"`
	RefreshTokenTTL string `mapstructure:"refresh_token_ttl"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() (*Config, error) {
	c := config.NewWithParams(&config.Params{
		Defaults: config.ConfigMap{
			"server.host":                "0.0.0.0",
			"server.port":                50001,
			"database.host":              "localhost",
			"database.port":              5432,
			"database.user":              "postgres",
			"database.password":          "postgres",
			"database.dbname":            "wms_auth",
			"database.sslmode":           "disable",
			"database.max_idle_conns":    10,
			"database.max_open_conns":    25,
			"database.conn_max_lifetime": "5m",
			"jwt.access_token_ttl":       "15m",
			"jwt.refresh_token_ttl":      "168h",
			"log.level":                  "debug",
		},
		EnvAlias: config.ConfigEnvAlias{
			"SERVER_HOST":          "server.host",
			"SERVER_PORT":          "server.port",
			"DB_HOST":              "database.host",
			"DB_PORT":              "database.port",
			"DB_USER":              "database.user",
			"DB_PASSWORD":          "database.password",
			"DB_NAME":              "database.dbname",
			"DB_SSLMODE":           "database.sslmode",
			"DB_MAX_IDLE_CONNS":    "database.max_idle_conns",
			"DB_MAX_OPEN_CONNS":    "database.max_open_conns",
			"DB_CONN_MAX_LIFETIME": "database.conn_max_lifetime",
			"JWT_PRIVATE_KEY_PATH": "jwt.private_key_path",
			"JWT_PUBLIC_KEY_PATH":  "jwt.public_key_path",
			"JWT_PRIVATE_KEY":      "jwt.private_key",
			"JWT_PUBLIC_KEY":       "jwt.public_key",
			"JWT_ACCESS_TTL":       "jwt.access_token_ttl",
			"JWT_REFRESH_TTL":      "jwt.refresh_token_ttl",
			"LOG_LEVEL":            "log.level",
		},
	})

	c.WithEnv()
	c.LoadConfigs()

	var cfg Config
	if err := c.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Handle JWT keys from files (ayotl can't read files into values)
	if cfg.JWT.PrivateKeyPath != "" && cfg.JWT.PrivateKeyPEM == "" {
		data, err := os.ReadFile(cfg.JWT.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}
		cfg.JWT.PrivateKeyPEM = string(data)
	}
	if cfg.JWT.PublicKeyPath != "" && cfg.JWT.PublicKeyPEM == "" {
		data, err := os.ReadFile(cfg.JWT.PublicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read public key: %w", err)
		}
		cfg.JWT.PublicKeyPEM = string(data)
	}

	if cfg.JWT.PrivateKeyPEM == "" || cfg.JWT.PublicKeyPEM == "" {
		return nil, fmt.Errorf("JWT private and public keys are required: set JWT_PRIVATE_KEY_PATH/JWT_PUBLIC_KEY_PATH or JWT_PRIVATE_KEY/JWT_PUBLIC_KEY env vars")
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

// AccessTokenTTLDuration parses AccessTokenTTL as a duration.
func (c *JWTConfig) AccessTokenTTLDuration() time.Duration {
	d, err := time.ParseDuration(c.AccessTokenTTL)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

// RefreshTokenTTLDuration parses RefreshTokenTTL as a duration.
func (c *JWTConfig) RefreshTokenTTLDuration() time.Duration {
	d, err := time.ParseDuration(c.RefreshTokenTTL)
	if err != nil {
		return 168 * time.Hour
	}
	return d
}

// Address returns the gRPC server address.
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LevelStr returns the slog-compatible log level.
func (c *LogConfig) LevelStr() string {
	return strings.ToLower(c.Level)
}
