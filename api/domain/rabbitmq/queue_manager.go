package rabbitmq

import (
	"log/slog"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type QueueManagerService struct {
	connection *RabbitMQConnectionService
}

func NewQueueManagerService(injector do.Injector) (*QueueManagerService, error) {
	conn := do.MustInvoke[*RabbitMQConnectionService](injector)

	return &QueueManagerService{
		connection: conn,
	}, nil
}

// DeclareQueue declares a queue with the given options
func (svc *QueueManagerService) DeclareQueue(queueName RoutingKey, options QueueOptions) error {
	channel := svc.connection.GetChannel()
	if channel == nil {
		return oops.New("channel is nil, connection may be closed")
	}

	_, err := channel.QueueDeclare(
		string(queueName),
		options.Durable,
		options.AutoDelete,
		options.Exclusive,
		options.NoWait,
		options.Args,
	)

	if err != nil {
		return oops.Wrapf(err, "failed to declare queue: %s", queueName)
	}

	slog.Info("Queue declared", "queue", queueName, "durable", options.Durable)
	return nil
}

// DeclareQueueForWorker declares a queue for a worker using its configuration
func (svc *QueueManagerService) DeclareQueueForWorker(worker Worker) error {
	queueName := worker.QueueName()
	options := worker.QueueOptions()

	// Use default options if QueueOptions returns zero value
	if options == nil {
		options = DefaultQueueOptions()
	}

	return svc.DeclareQueue(queueName, *options)
}

// EnsureQueueExists checks if a queue exists and creates it if it doesn't
// This is idempotent - safe to call multiple times
func (svc *QueueManagerService) EnsureQueueExists(queueName RoutingKey, options QueueOptions) error {
	return svc.DeclareQueue(queueName, options)
}
