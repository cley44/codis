package rabbitmq

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type PublisherService struct {
	connection *RabbitMQConnectionService
}

func NewPublisherService(injector do.Injector) (*PublisherService, error) {
	conn := do.MustInvoke[*RabbitMQConnectionService](injector)
	return &PublisherService{connection: conn}, nil
}

func (s *PublisherService) Publish(routingKey RoutingKey, msgBody AMQPMessageBody) error {
	channel := s.connection.GetChannel()
	if channel == nil {
		return oops.New("channel is nil, connection may be closed")
	}

	msgBody.EmittedAt = time.Now()
	msgBody.Type = routingKey

	body, err := serialize(msgBody)
	if err != nil {
		return oops.Wrapf(err, "failed to marshal message")
	}

	return channel.PublishWithContext(context.Background(), "", string(routingKey), false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
	})
}
