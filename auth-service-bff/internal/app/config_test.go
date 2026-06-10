package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDefaults(t *testing.T) {
	// No env vars set — should use defaults
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8081, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.AuthService.GRPCHost)
	assert.Equal(t, 50001, cfg.AuthService.GRPCPort)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "*", cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"*"}, cfg.CORS.Origins())
}

func TestConfigCustom(t *testing.T) {
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("AUTH_SERVICE_GRPC_HOST", "login.internal")
	t.Setenv("AUTH_SERVICE_GRPC_PORT", "50002")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://app.example.com")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "login.internal", cfg.AuthService.GRPCHost)
	assert.Equal(t, 50002, cfg.AuthService.GRPCPort)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "http://localhost:3000,https://app.example.com", cfg.CORS.AllowedOrigins)
	origins := cfg.CORS.Origins()
	assert.Len(t, origins, 2)
	assert.Equal(t, "http://localhost:3000", origins[0])
	assert.Equal(t, "https://app.example.com", origins[1])
}

func TestCORSOriginsEmpty(t *testing.T) {
	cfg := &CORSConfig{}
	assert.Equal(t, []string{"*"}, cfg.Origins())
}

func TestCORSOriginsCustom(t *testing.T) {
	cfg := &CORSConfig{
		AllowedOrigins: "http://example.com",
	}
	assert.Equal(t, []string{"http://example.com"}, cfg.Origins())
}
