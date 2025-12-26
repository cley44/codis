package rabbitmq

import (
	"codis/config"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type RabbitMQConnectionService struct {
	config      *config.ConfigService
	conn        *amqp091.Connection
	channel     *amqp091.Channel
	isConnected bool
}

func NewRabbitMQConnectionService(injector do.Injector) (*RabbitMQConnectionService, error) {
	cfg := do.MustInvoke[*config.ConfigService](injector)

	vhost := cfg.RabbitMQ.VHost
	if vhost == "" {
		vhost = "/"
	}

	uri := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Hostname,
		cfg.RabbitMQ.Port,
		vhost,
	)

	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to connect to RabbitMQ")
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, oops.Wrapf(err, "failed to open channel")
	}

	service := &RabbitMQConnectionService{
		config:      cfg,
		conn:        conn,
		channel:     channel,
		isConnected: true,
	}

	slog.Info("Connected to RabbitMQ", "host", cfg.RabbitMQ.Hostname, "port", cfg.RabbitMQ.Port)

	return service, nil
}

func (svc *RabbitMQConnectionService) GetChannel() *amqp091.Channel {
	return svc.channel
}

func (svc *RabbitMQConnectionService) GetConnection() *amqp091.Connection {
	return svc.conn
}

func (svc *RabbitMQConnectionService) IsConnected() bool {
	return svc.isConnected && !svc.conn.IsClosed()
}

func (svc *RabbitMQConnectionService) Reconnect() error {
	if svc.conn != nil && !svc.conn.IsClosed() {
		svc.conn.Close()
	}

	vhost := svc.config.RabbitMQ.VHost
	if vhost == "" {
		vhost = "/"
	}

	uri := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		svc.config.RabbitMQ.Username,
		svc.config.RabbitMQ.Password,
		svc.config.RabbitMQ.Hostname,
		svc.config.RabbitMQ.Port,
		vhost,
	)

	conn, err := amqp091.Dial(uri)
	if err != nil {
		svc.isConnected = false
		return oops.Wrapf(err, "failed to reconnect to RabbitMQ")
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		svc.isConnected = false
		return oops.Wrapf(err, "failed to open channel after reconnect")
	}

	svc.conn = conn
	svc.channel = channel
	svc.isConnected = true

	slog.Info("Reconnected to RabbitMQ")
	return nil
}

func (svc *RabbitMQConnectionService) Close() error {
	var err error
	if svc.channel != nil {
		if closeErr := svc.channel.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if svc.conn != nil {
		if closeErr := svc.conn.Close(); closeErr != nil {
			if err != nil {
				err = oops.Wrapf(closeErr, "failed to close connection")
			} else {
				err = closeErr
			}
		}
	}
	svc.isConnected = false
	return err
}
