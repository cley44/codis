package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type ConsumerManagerService struct {
	connection   *RabbitMQConnectionService
	queueManager *QueueManagerService
	workers      []Worker
	consumers    map[RoutingKey]*consumerContext
	mu           sync.RWMutex
}

type consumerContext struct {
	cancel    context.CancelFunc
	done      chan struct{}
	queueName string
}

func NewConsumerManagerService(injector do.Injector) (*ConsumerManagerService, error) {
	conn := do.MustInvoke[*RabbitMQConnectionService](injector)
	queueMgr := do.MustInvoke[*QueueManagerService](injector)

	return &ConsumerManagerService{
		connection:   conn,
		queueManager: queueMgr,
		workers:      make([]Worker, 0),
		consumers:    make(map[RoutingKey]*consumerContext),
	}, nil
}

// RegisterWorker registers a worker to be started later
func (svc *ConsumerManagerService) RegisterWorker(worker Worker) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.workers = append(svc.workers, worker)
	slog.Info("Worker registered", "queue", worker.QueueName())
}

// StartAll starts all registered workers
// This should be called after queues are declared
func (svc *ConsumerManagerService) StartAll() error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	for _, worker := range svc.workers {
		if err := svc.startWorker(worker); err != nil {
			return oops.Wrapf(err, "failed to start worker for queue: %s", worker.QueueName())
		}
	}

	slog.Info("All workers started", "count", len(svc.workers))
	return nil
}

// startWorker starts a single worker in a goroutine
func (svc *ConsumerManagerService) startWorker(worker Worker) error {
	queueName := worker.QueueName()

	// Check if already running
	if _, exists := svc.consumers[queueName]; exists {
		slog.Warn("Worker already running", "queue", queueName)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	svc.consumers[queueName] = &consumerContext{
		cancel:    cancel,
		done:      done,
		queueName: string(queueName),
	}

	// Start consumer in goroutine
	go svc.consumeMessages(ctx, done, worker)

	slog.Info("Worker started", "queue", queueName)
	return nil
}

// consumeMessages handles message consumption for a worker
func (svc *ConsumerManagerService) consumeMessages(ctx context.Context, done chan struct{}, worker Worker) {
	defer close(done)

	queueName := worker.QueueName()
	channel := svc.connection.GetChannel()

	if channel == nil {
		slog.Error("Channel is nil, cannot consume", "queue", queueName)
		return
	}

	msgs, err := channel.Consume(
		string(queueName),
		"",    // consumer tag (empty = auto-generated)
		false, // auto-ack (false = manual ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		slog.Error("Failed to register consumer", "queue", queueName, "error", oops.Wrap(err))
		return
	}

	slog.Info("Consumer registered, waiting for messages", "queue", queueName)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Consumer stopped", "queue", queueName)
			return
		case msg, ok := <-msgs:
			if !ok {
				slog.Warn("Message channel closed", "queue", queueName)
				return
			}

			// Adapt and unmarshal msg to our own message type
			body := AMQPMessageBody{}
			if err := json.Unmarshal(msg.Body, &body); err != nil {
				slog.Error("Failed to unmarshal message body", "queue", queueName, "error", oops.Wrap(err))
				return
			}

			amqpMsg := AMQPMessage{
				RoutingKey: RoutingKey(msg.RoutingKey),
				Delivery:   msg,
				Body:       body,
			}

			// Process message
			if err := worker.HandleMessage(amqpMsg); err != nil {
				slog.Error("Failed to process message", "queue", queueName, "error", oops.Wrap(err))
				// Reject message and requeue it
				if nackErr := msg.Nack(false, true); nackErr != nil {
					slog.Error("Failed to nack message", "queue", queueName, "error", nackErr)
				}
			} else {
				// Acknowledge message
				if ackErr := msg.Ack(false); ackErr != nil {
					slog.Error("Failed to ack message", "queue", queueName, "error", ackErr)
				}
			}
		}
	}
}

// StopAll gracefully stops all running consumers
func (svc *ConsumerManagerService) StopAll() error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	slog.Info("Stopping all consumers", "count", len(svc.consumers))

	for queueName, ctx := range svc.consumers {
		slog.Info("Stopping consumer", "queue", queueName)
		ctx.cancel()
	}

	// Wait for all consumers to finish
	for queueName, ctx := range svc.consumers {
		<-ctx.done
		slog.Info("Consumer stopped", "queue", queueName)
	}

	svc.consumers = make(map[RoutingKey]*consumerContext)
	return nil
}

// StopWorker stops a specific worker by queue name
func (svc *ConsumerManagerService) StopWorker(queueName RoutingKey) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	ctx, exists := svc.consumers[queueName]
	if !exists {
		return oops.New("worker not found")
	}

	ctx.cancel()
	<-ctx.done
	delete(svc.consumers, queueName)

	slog.Info("Worker stopped", "queue", queueName)
	return nil
}

// GetRegisteredWorkers returns a list of all registered worker queue names
func (svc *ConsumerManagerService) GetRegisteredWorkers() []RoutingKey {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	queues := make([]RoutingKey, len(svc.workers))
	for i, worker := range svc.workers {
		queues[i] = worker.QueueName()
	}
	return queues
}

// DeclareAllQueues declares queues for all registered workers
func (svc *ConsumerManagerService) DeclareAllQueues() error {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	for _, worker := range svc.workers {
		if err := svc.queueManager.DeclareQueueForWorker(worker); err != nil {
			return oops.Wrapf(err, "failed to declare queue for worker: %s", worker.QueueName())
		}
	}

	return nil
}
