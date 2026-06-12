package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDefaults(t *testing.T) {
	// No env vars set — should fail because JWT keys are required
	_, err := LoadConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "JWT private and public keys are required")
}

func TestConfigWithEnv(t *testing.T) {
	t.Setenv("JWT_PRIVATE_KEY", "test-private-key")
	t.Setenv("JWT_PUBLIC_KEY", "test-public-key")
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "admin")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("DB_MAX_IDLE_CONNS", "5")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_CONN_MAX_LIFETIME", "10m")
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("JWT_ACCESS_TTL", "30m")
	t.Setenv("JWT_REFRESH_TTL", "48h")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "127.0.0.1:9090", cfg.Server.Address())

	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
	assert.Equal(t, "admin", cfg.Database.User)
	assert.Equal(t, "secret", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.DBName)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10*time.Minute, cfg.Database.ConnMaxLifetimeDuration())

	assert.Equal(t, "warn", cfg.Log.LevelStr())

	assert.Equal(t, 30*time.Minute, cfg.JWT.AccessTokenTTLDuration())
	assert.Equal(t, 48*time.Hour, cfg.JWT.RefreshTokenTTLDuration())
	assert.Equal(t, "test-private-key", cfg.JWT.PrivateKeyPEM)
	assert.Equal(t, "test-public-key", cfg.JWT.PublicKeyPEM)
}

func TestConfigDSN(t *testing.T) {
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

func TestServerConfigAddress(t *testing.T) {
	cfg := &ServerConfig{Host: "0.0.0.0", Port: 50001}
	assert.Equal(t, "0.0.0.0:50001", cfg.Address())
}

func TestLogConfigLevelStr(t *testing.T) {
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

func TestConfigMissingPrivateKey(t *testing.T) {
	t.Setenv("JWT_PUBLIC_KEY", "pub-key")
	// Private key missing
	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "JWT private and public keys are required")
}

func TestConfigMissingPublicKey(t *testing.T) {
	t.Setenv("JWT_PRIVATE_KEY", "priv-key")
	// Public key missing
	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "JWT private and public keys are required")
}
