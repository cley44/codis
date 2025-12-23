package rabbitmq

// QueueOptions defines configuration options for a RabbitMQ queue
type QueueOptions struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]interface{}
}

// DefaultQueueOptions returns sensible default queue options
func DefaultQueueOptions() *QueueOptions {
	return &QueueOptions{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

// Worker interface defines the contract for RabbitMQ message workers
type Worker interface {
	// QueueName returns the name of the queue this worker consumes from
	QueueName() RoutingKey

	// HandleMessage processes a single message from the queue
	// Returns an error if message processing failed
	HandleMessage(msg AMQPMessage) error

	// QueueOptions returns the queue configuration options
	// If nil is returned, DefaultQueueOptions() will be used
	QueueOptions() *QueueOptions
}
