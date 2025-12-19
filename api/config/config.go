package config

import (
	"github.com/samber/config"
	"github.com/samber/do/v2"
)

type ConfigService struct {
	Discord         DiscordConfig
	Postgres        PostgresConfig
	Auth            AuthConfig
	Instrumentation InstrumentationConfig
}

func NewConfigService(injector do.Injector) (*ConfigService, error) {
	c := ConfigService{}

	// @TODO change this when using docker
	// config.FromEnv().To(&c)
	config.From(".env").To(&c)

	return &c, nil
}
