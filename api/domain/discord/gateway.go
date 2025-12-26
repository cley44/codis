package discord

import (
	"codis/domain/rabbitmq"
	"codis/models"
	"codis/utils/slogger"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

func (m DiscordService) OnEvent(event bot.Event) {
	msg := rabbitmq.AMQPMessageBody{
		DiscordEvent: rabbitmq.DiscordEvent{},
	}
	switch e := event.(type) {
	case *events.MessageCreate:
		msg.DiscordEvent.Type = models.DiscordEventTypeMessageCreate
		msg.DiscordEvent.GuildID = e.GuildID.String()
		msg.DiscordEvent.UserID = e.Message.Author.ID.String()
	case *events.MessageReactionAdd:
	case *events.GuildMessageReactionAdd:
		msg.DiscordEvent.Type = models.DiscordEventTypeMessageReactionAdd
		msg.DiscordEvent.GuildID = e.GuildID.String()
		msg.DiscordEvent.UserID = e.UserID.String()
	case *events.Ready:
	case *events.HeartbeatAck:
		// Ignored events
		return
	default:
		slogger.Info("Unknown event type", "event", event)
		return
	}

	// msg.DiscordEvent.RawEvent = event

	err := m.publisherService.Publish(rabbitmq.RoutingKeyDispatch, msg)
	if err != nil {
		slogger.Errorf("Failed to publish message", "error", err)
	}
}
