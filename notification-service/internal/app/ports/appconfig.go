// Package ports holds shared interfaces and config types to avoid import cycles.
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
	Sender struct {
		Type string `yaml:"type"`
	} `yaml:"sender"`
	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
	Auth struct {
		APIKey        string   `yaml:"api_key"`
		ExemptMethods []string `yaml:"exempt_methods"`
	} `yaml:"auth"`
}

// AppConfig is the interface for loading and accessing configuration.
type AppConfig interface {
	LoadConfigs() error
	GetConfigs() *Config
}
