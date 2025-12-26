package rabbitmq

import (
	"codis/models"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RoutingKey string

const (
	RoutingKeyExample     RoutingKey = "example"
	RoutingKeyDispatch    RoutingKey = "dispatch"
	RoutingKeyNodeExecute RoutingKey = "node_execute"
)

type AMQPMessageBody struct {
	Type      RoutingKey `amqp:"type"`
	EmittedAt time.Time  `amqp:"emitted_at"`

	DiscordEvent DiscordEvent `amqp:"discord_event"`
}

type DiscordEvent struct {
	Type            models.DiscordEventType `amqp:"discord_event_type"`
	NodeIDToExecute string                  `amqp:"node_id_to_execute"`
	GuildID         string                  `amqp:"guild_id"`
	UserID          string                  `amqp:"user_id"`
	RoleID          string                  `amqp:"role_id"`

	// RawEvent bot.Event `amqp:"raw_event"`
}

type AMQPMessage struct {
	RoutingKey RoutingKey
	Delivery   amqp091.Delivery
	Body       AMQPMessageBody
}

type NodeHandler interface {
	GetType() models.DiscordNodeType
	Execute(msg AMQPMessageBody) error
}
