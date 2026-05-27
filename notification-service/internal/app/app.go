package app

import (
	"go.uber.org/zap"

	"github.com/beabys/wms/notification-service/internal/app/ports"
)

// New creates a new App instance.
func New() *App {
	return &App{}
}

// SetConfigs sets the app configuration and loads it.
func (app *App) SetConfigs(cfg ports.AppConfig) error {
	app.Config = cfg
	return cfg.LoadConfigs()
}

// Setup initializes the logger.
func (app *App) Setup(configs ports.AppConfig) error {
	if err := app.SetConfigs(configs); err != nil {
		return err
	}

	appConfig := configs.GetConfigs()

	var logger *zap.Logger
	var err error
	if appConfig.Service.Env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	app.Logger = logger

	return nil
}

// GetLogger returns the app logger.
func (app *App) GetLogger() *zap.Logger {
	return app.Logger
}
