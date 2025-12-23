package discord

import (
	"codis/config"
	"codis/domain/rabbitmq"
	"codis/models"
	"codis/repository"
	"codis/utils/slogger"
	"context"
	"net/http"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"

	"github.com/samber/do/v2"
)

type DiscordService struct {
	config           *config.ConfigService
	oauthClient      oauth2.Client
	botClient        bot.Client
	userRepository   *repository.UserRepository
	publisherService *rabbitmq.PublisherService
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

	m := DiscordService{
		config:           config,
		oauthClient:      oauthClient,
		botClient:        nil,
		userRepository:   do.MustInvoke[*repository.UserRepository](injector),
		publisherService: do.MustInvoke[*rabbitmq.PublisherService](injector),
	}

	botClient, err := disgo.New(config.Discord.DiscordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				// gateway.IntentGuildMembers,
				// gateway.IntentGuilds,
				gateway.IntentsNonPrivileged,
			),
			gateway.WithAutoReconnect(true),
		),
		bot.WithEventListenerFunc(m.OnEvent),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
		),
		// bot.WithLogger(slog.Default()),
		// bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		panic(err)
	}

	m.botClient = botClient

	botClient.OpenGateway(context.Background())

	return &m, nil
}

func (m DiscordService) OnEvent(event bot.Event) {
	msg := rabbitmq.AMQPMessageBody{
		DiscordEvent: rabbitmq.DiscordEvent{},
	}
	switch e := event.(type) {
	case *events.MessageCreate:
		msg.DiscordEvent.MessageCreateEvent = e
		msg.DiscordEvent.Type = models.DiscordEventTypeMessageCreate
	case *events.MessageReactionAdd:
		msg.DiscordEvent.MessageReactionAddEvent = e
		msg.DiscordEvent.Type = models.DiscordEventTypeMessageReactionAdd
	case *events.Ready:
	case *events.HeartbeatAck:
		// Ignored events
	default:
		slogger.Info("Unknown event type", "event", event)
		return
	}

	m.publisherService.Publish(rabbitmq.RoutingKeyDispatch, msg)

}
