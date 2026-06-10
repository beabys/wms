package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomerBFFConfigDefaults(t *testing.T) {
	// No env vars set — should use defaults
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8082, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.CustomerService.GRPCHost)
	assert.Equal(t, 50002, cfg.CustomerService.GRPCPort)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "*", cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"*"}, cfg.CORS.Origins())
}

func TestCustomerBFFConfigCustom(t *testing.T) {
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("CUSTOMER_SERVICE_GRPC_HOST", "customer.internal")
	t.Setenv("CUSTOMER_SERVICE_GRPC_PORT", "50003")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://app.example.com")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "customer.internal", cfg.CustomerService.GRPCHost)
	assert.Equal(t, 50003, cfg.CustomerService.GRPCPort)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "http://localhost:3000,https://app.example.com", cfg.CORS.AllowedOrigins)
	origins := cfg.CORS.Origins()
	assert.Len(t, origins, 2)
	assert.Equal(t, "http://localhost:3000", origins[0])
	assert.Equal(t, "https://app.example.com", origins[1])
}

func TestCustomerBFFCORSOriginsEmpty(t *testing.T) {
	cfg := &CORSConfig{}
	assert.Equal(t, []string{"*"}, cfg.Origins())
}

func TestCustomerBFFCORSOriginsCustom(t *testing.T) {
	cfg := &CORSConfig{
		AllowedOrigins: "http://example.com",
	}
	assert.Equal(t, []string{"http://example.com"}, cfg.Origins())
}

func TestCustomerBFFServerAddress(t *testing.T) {
	cfg := &ServerConfig{Host: "0.0.0.0", Port: 8082}
	assert.Equal(t, "0.0.0.0:8082", cfg.Address())
}

func TestCustomerBFFLevelStr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty default", "", "debug"},
		{"debug", "debug", "debug"},
		{"info", "info", "info"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LogConfig{Level: tt.input}
			assert.Equal(t, tt.want, c.LevelStr())
		})
	}
}
