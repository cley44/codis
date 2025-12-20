package discord

import (
	"codis/config"
	"codis/repository"
	"context"
	"net/http"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"

	"github.com/samber/do/v2"
)

type DiscordService struct {
	config         *config.ConfigService
	oauthClient    oauth2.Client
	botClient      bot.Client
	userRepository *repository.UserRepository
}

func NewDiscordService(injector do.Injector) (*DiscordService, error) {
	config := do.MustInvoke[*config.ConfigService](injector)

	clientID, err := snowflake.Parse(config.Discord.ClientID)
	if err != nil {
		panic("Discord Client ID is not defined")
	}

	// slog.SetLogLoggerLevel(slog.LevelDebug)

	oauthClient := oauth2.New(
		clientID,
		config.Discord.ClientSecret,
		oauth2.WithRestClientConfigOpts(rest.WithHTTPClient(http.DefaultClient)),
		// oauth2.WithLogger(slog.Default()),
	)

	botClient, err := disgo.New(config.Discord.DiscordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMembers,
				gateway.IntentGuilds,
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
		),
		// bot.WithLogger(slog.Default()),
		// bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		panic(err)
	}

	botClient.OpenGateway(context.Background())

	m := DiscordService{
		config:         config,
		oauthClient:    oauthClient,
		botClient:      botClient,
		userRepository: do.MustInvoke[*repository.UserRepository](injector),
	}

	return &m, nil
}
