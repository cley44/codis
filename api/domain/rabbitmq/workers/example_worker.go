package workers

import (
	"codis/domain/rabbitmq"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

// ExampleWorker demonstrates how to implement a RabbitMQ worker
// This is a reference implementation that you can copy and modify for your own workers
type ExampleWorker struct {
	// Add your dependencies here (e.g., repositories, services, etc.)
}

// NewExampleWorker creates a new example worker instance
func NewExampleWorker() *ExampleWorker {
	return &ExampleWorker{}
}

// QueueName returns the name of the queue this worker consumes from
func (w *ExampleWorker) QueueName() rabbitmq.RoutingKey {
	return rabbitmq.RoutingKeyExample
}

// QueueOptions returns the queue configuration options
func (w *ExampleWorker) QueueOptions() *rabbitmq.QueueOptions {
	return &rabbitmq.QueueOptions{
		Durable:    true,  // Queue survives broker restart
		AutoDelete: false, // Queue is not deleted when no longer used
		Exclusive:  false, // Queue can be accessed by other connections
		NoWait:     false, // Wait for server confirmation
		Args:       nil,   // Additional queue arguments
	}
}

// HandleMessage processes a single message from the queue
func (w *ExampleWorker) HandleMessage(msg rabbitmq.AMQPMessage) error {
	slog.Info("Processing message",
		"queue", w.QueueName(),
		"routing_key", msg.RoutingKey,
	)

	// Example: Process the message
	slog.Info("Message payload", "payload", msg.Body)

	// Example: Do some work
	// result := w.processData(payload)

	// If processing fails, return an error
	// The message will be nacked and requeued
	// If processing succeeds, return nil
	// The message will be acked

	return nil
}

// Example of a worker with custom queue options
type CustomOptionsWorker struct {
	queueName string
}

func NewCustomOptionsWorker(queueName string) *CustomOptionsWorker {
	return &CustomOptionsWorker{
		queueName: queueName,
	}
}

func (w *CustomOptionsWorker) QueueName() string {
	return w.queueName
}

func (w *CustomOptionsWorker) QueueOptions() rabbitmq.QueueOptions {
	return rabbitmq.QueueOptions{
		Durable:    true,
		AutoDelete: true, // Queue will be deleted when no longer used
		Exclusive:  false,
		NoWait:     false,
		Args: map[string]interface{}{
			"x-message-ttl": 60000, // Messages expire after 60 seconds
		},
	}
}

func (w *CustomOptionsWorker) HandleMessage(msg amqp091.Delivery) error {
	slog.Info("Custom worker processing message",
		"queue", w.queueName,
		"body", string(msg.Body),
	)

	// Process message here
	return nil
}
