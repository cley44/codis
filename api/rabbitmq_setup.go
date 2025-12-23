package codis

import (
	rabbitmqDomain "codis/domain/rabbitmq"
	"log/slog"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

// SetupRabbitMQWorkers declares queues and starts consumers for all registered workers
// This should be called after workers are registered but before starting the HTTP server
func SetupRabbitMQWorkers(injector do.Injector) error {
	consumerManager := do.MustInvoke[*rabbitmqDomain.ConsumerManagerService](injector)

	// Declare all queues for registered workers
	if err := consumerManager.DeclareAllQueues(); err != nil {
		return oops.Wrapf(err, "failed to declare queues")
	}

	// Start all consumers
	if err := consumerManager.StartAll(); err != nil {
		return oops.Wrapf(err, "failed to start consumers")
	}

	workerCount := len(consumerManager.GetRegisteredWorkers())
	slog.Info("RabbitMQ workers setup completed", "worker_count", workerCount)
	return nil
}

// RegisterWorkers is a helper function to register workers with the consumer manager
// This should be called before SetupRabbitMQWorkers
func RegisterWorkers(injector do.Injector, workers ...rabbitmqDomain.Worker) {
	consumerManager := do.MustInvoke[*rabbitmqDomain.ConsumerManagerService](injector)
	for _, worker := range workers {
		consumerManager.RegisterWorker(worker)
	}
}
