package app

import (
	grpcdapter "github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/grpc"
	httpadapter "github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/http"
	"github.com/beabys/wms/pkg/logger"
)

// App is the application struct holding all dependencies.
type App struct {
	Config     *Config
	Logger     logger.Logger
	HTTPServer *httpadapter.HTTPServer
	AuthClient *grpcdapter.AuthClient
}

// Config holds application configuration.
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	AuthService AuthServiceConfig `mapstructure:"auth_service"`
	Log         LogConfig         `mapstructure:"log"`
	CORS        CORSConfig        `mapstructure:"cors"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// AuthServiceConfig holds gRPC connection configuration for the auth service.
type AuthServiceConfig struct {
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
