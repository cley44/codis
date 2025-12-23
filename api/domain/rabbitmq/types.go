package rabbitmq

import (
	"codis/models"
	"time"

	"github.com/disgoorg/disgo/events"
	"github.com/rabbitmq/amqp091-go"
)

type RoutingKey string

const (
	RoutingKeyExample  RoutingKey = "example"
	RoutingKeyDispatch RoutingKey = "dispatch"
)

type AMQPMessageBody struct {
	Type      RoutingKey `amqp:"type"`
	EmittedAt time.Time  `amqp:"emitted_at"`

	DiscordEvent DiscordEvent `amqp:"discord_event"`
}

type DiscordEvent struct {
	Type                    models.DiscordEventType    `amqp:"discord_event_type"`
	MessageCreateEvent      *events.MessageCreate      `amqp:"message_create_event"`
	MessageReactionAddEvent *events.MessageReactionAdd `amqp:"message_reaction_add_event"`
}

type AMQPMessage struct {
	RoutingKey RoutingKey
	Delivery   amqp091.Delivery
	Body       AMQPMessageBody
}
