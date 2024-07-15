package config

import (
	"fmt"

	"github.com/samber/config"
	"github.com/samber/do/v2"
)

type ConfigService struct {
	Discord  DiscordConfig
	Postgres PostgresConfig
}

func NewConfigService(injector do.Injector) (*ConfigService, error) {
	c := ConfigService{}

	// @TODO change this when using docker
	// config.FromEnv().To(&c)
	config.From(".env").To(&c)

	// Computed env vars can be added here
	c.Postgres.URI = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.Postgres.Username, c.Postgres.Password, c.Postgres.Hostname, c.Postgres.Port, c.Postgres.Database)

	return &c, nil
}
