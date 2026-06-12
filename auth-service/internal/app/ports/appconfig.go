package ports

type AppConfig interface {
	LoadConfigs() error
	GetConfigs() interface{}
}
