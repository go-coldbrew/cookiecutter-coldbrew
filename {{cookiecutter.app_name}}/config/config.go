package config

import (
	"context"
	"log/slog"

	cbConfig "github.com/go-coldbrew/core/config"
	"github.com/kelseyhightower/envconfig"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service/auth"
)

// defaultConfig is the default configuration for the application
// It is loaded from environment variables
var defaultConfig Config

type Config struct {
	cbConfig.Config
	auth.AuthConfig
	PanicOnConfigError bool `envconfig:"PANIC_ON_CONFIG_ERROR" default:"true"`
	// App configuration
	// Remove this line and add your own configuration
	Prefix string `envconfig:"PREFIX" default:"got"`
}

func init() {
	err := envconfig.Process("", &defaultConfig)
	// fail on error
	if err != nil {
		if defaultConfig.PanicOnConfigError {
			panic(err)
		} else {
			slog.LogAttrs(context.Background(), slog.LevelError, "error while loading config", slog.Any("err", err))
		}
	}
}

// Get returns the default configuration
// This is used by the application to load the configuration
func Get() Config {
	return defaultConfig
}

// GetColdBrewConfig returns the default configuration
// This is used by the coldbrew framework to load the configuration for the application
func GetColdBrewConfig() cbConfig.Config {
	return defaultConfig.Config
}
