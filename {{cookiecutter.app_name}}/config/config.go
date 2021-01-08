package config

import (
	cbConfig "github.com/go-coldbrew/core/config"
	"github.com/kelseyhightower/envconfig"
)

var defaultConfig Config

type Config struct {
	cbConfig.Config
	// App configuration
	Prefix string `envconfig:"PREFIX" default:"got"`
}

func init() {
	envconfig.Process("", &defaultConfig)
	// fail on error
}

func Get() Config {
	return defaultConfig
}

func GetColdBrewConfig() cbConfig.Config {
	return defaultConfig.Config
}
