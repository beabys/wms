package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomerConfigDefaults(t *testing.T) {
	// No env vars set — should fail because JWT public key is required
	_, err := LoadConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "JWT public key is required")
}

func TestAuthServiceConfigAddress(t *testing.T) {
	cfg := &AuthServiceConfig{Host: "auth.example.com", Port: 50001}
	assert.Equal(t, "auth.example.com:50001", cfg.Address())
}

func TestAuthServiceConfigDefaults(t *testing.T) {
	cfg := &AuthServiceConfig{}
	assert.Equal(t, ":0", cfg.Address())
}

func TestCustomerConfigWithEnv(t *testing.T) {
	t.Setenv("JWT_PUBLIC_KEY", "test-public-key")
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5444")
	t.Setenv("DB_USER", "admin")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "customers")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("DB_MAX_IDLE_CONNS", "5")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_CONN_MAX_LIFETIME", "10m")
	t.Setenv("LOG_LEVEL", "warn")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "127.0.0.1:9090", cfg.Server.Address())

	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, 5444, cfg.Database.Port)
	assert.Equal(t, "admin", cfg.Database.User)
	assert.Equal(t, "secret", cfg.Database.Password)
	assert.Equal(t, "customers", cfg.Database.DBName)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10*time.Minute, cfg.Database.ConnMaxLifetimeDuration())

	assert.Equal(t, "warn", cfg.Log.LevelStr())

	assert.Equal(t, "test-public-key", cfg.JWT.PublicKeyPEM)
}

func TestCustomerConfigDSN(t *testing.T) {
	cfg := &DatabaseConfig{
		User:     "user",
		Password: "pass",
		Host:     "host",
		Port:     1234,
		DBName:   "db",
		SSLMode:  "disable",
	}
	expected := "postgres://user:pass@host:1234/db?sslmode=disable"
	assert.Equal(t, expected, cfg.DSN())
}

func TestCustomerServerConfigAddress(t *testing.T) {
	cfg := &ServerConfig{Host: "0.0.0.0", Port: 50002}
	assert.Equal(t, "0.0.0.0:50002", cfg.Address())
}

func TestCustomerLogConfigLevelStr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "debug", "debug"},
		{"uppercase", "INFO", "info"},
		{"mixed", "Warn", "warn"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LogConfig{Level: tt.input}
			assert.Equal(t, tt.want, c.LevelStr())
		})
	}
}

func TestCustomerConnMaxLifetimeDurationDefault(t *testing.T) {
	cfg := &DatabaseConfig{ConnMaxLifetime: ""}
	assert.Equal(t, 5*time.Minute, cfg.ConnMaxLifetimeDuration())
}

func TestCustomerConnMaxLifetimeDurationCustom(t *testing.T) {
	cfg := &DatabaseConfig{ConnMaxLifetime: "30s"}
	assert.Equal(t, 30*time.Second, cfg.ConnMaxLifetimeDuration())
}

func TestCustomerConfigMissingPublicKey(t *testing.T) {
	// No JWT keys set
	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "JWT public key is required")
}
