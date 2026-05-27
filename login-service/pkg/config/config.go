package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the login-service configuration.
type Config struct {
	DBDsn             string `yaml:"db_dsn"`
	JWTPrivateKeyPath string `yaml:"jwt_private_key_path"`
	JWTPublicKeyPath  string `yaml:"jwt_public_key_path"`
	AccessTokenTTL    string `yaml:"access_token_ttl"`
	RefreshTokenTTL   string `yaml:"refresh_token_ttl"`
	GRPCPort          int    `yaml:"grpc_port"`
}

// New returns a Config with default values.
func New() *Config {
	return &Config{
		AccessTokenTTL:  "15m",
		RefreshTokenTTL: "168h",
		GRPCPort:        50001,
	}
}

// LoadConfigs loads configuration from a YAML file.
func (c *Config) LoadConfigs() error {
	cfgPath := "login-service/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		cfgPath = envPath
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

// GetConfigs returns the config pointer.
func (c *Config) GetConfigs() *Config {
	return c
}

// ParsePrivateKey loads an RSA private key from a PEM file.
func ParsePrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM data in private key file")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA private key")
	}

	return rsaKey, nil
}

// ReadPublicKeyPEM reads the raw PEM bytes from a public key file.
func ReadPublicKeyPEM(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read public key PEM file: %w", err)
	}
	return string(data), nil
}

// ParsePublicKey loads an RSA public key from a PEM file.
func ParsePublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read public key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM data in public key file")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA public key")
	}

	return rsaKey, nil
}
