package discord

import (
	"codis/config"
	"codis/repository"
	"log/slog"
	"net/http"

	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"

	"github.com/samber/do/v2"
)

type DiscordService struct {
	config         *config.ConfigService
	oauthClient    oauth2.Client
	userRepository *repository.UserRepository
}

func NewDiscordService(injector do.Injector) (*DiscordService, error) {
	config := do.MustInvoke[*config.ConfigService](injector)

	clientID, err := snowflake.Parse(config.Discord.ClientID)
	if err != nil {
		panic("Discord Client ID is not defined")
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	client := oauth2.New(
		clientID,
		config.Discord.ClientSecret,
		oauth2.WithRestClientConfigOpts(rest.WithHTTPClient(http.DefaultClient)),
		oauth2.WithLogger(slog.Default()),
	)

	m := DiscordService{
		config:         config,
		oauthClient:    client,
		userRepository: do.MustInvoke[*repository.UserRepository](injector),
	}

	return &m, nil
}
