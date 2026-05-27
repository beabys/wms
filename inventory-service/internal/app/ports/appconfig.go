// Package ports holds shared interfaces to avoid import cycles
// when interfaces reference types from their own package.
//
// The pattern: if interface A in package P uses type B from package P
// in its method signature, generating mocks for A requires importing P,
// which creates a cycle when tests in P import those mocks.
//
// Moving such interfaces to a sibling ports/ package breaks the cycle
// while keeping the consumer-defines-interface principle intact.
package ports

// Config holds application configuration loaded from config.yaml.
type Config struct {
	Service struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
	} `yaml:"service"`
	GRPC struct {
		Port int `yaml:"port"`
	} `yaml:"grpc"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`
	Auth struct {
		JWKSURL       string   `yaml:"jwks_url"`
		ExemptMethods []string `yaml:"exempt_methods"`
	} `yaml:"auth"`

	JWTPublicKeyPath string `yaml:"jwt_public_key_path"`
}

// AppConfig is the interface for loading and accessing configuration.
type AppConfig interface {
	LoadConfigs() error
	GetConfigs() *Config
}
