package main

import (
	"codis"
	rabbitmqDomain "codis/domain/rabbitmq"
	"codis/domain/rabbitmq/workers"
	"log/slog"
	"syscall"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

func main() {
	injector := codis.RegisterAll()

	// Register RabbitMQ workers here
	// Example:
	codis.RegisterWorkers(injector, workers.NewDispatchWorker(injector))

	// Setup RabbitMQ: declare queues and start consumers
	if err := codis.SetupRabbitMQWorkers(injector); err != nil {
		slog.Error("Failed to setup RabbitMQ workers", "error", oops.Wrap(err))
		// Continue anyway - workers are optional
	}

	// Get consumer manager for graceful shutdown
	consumerManager := do.MustInvoke[*rabbitmqDomain.ConsumerManagerService](injector)
	connection := do.MustInvoke[*rabbitmqDomain.RabbitMQConnectionService](injector)

	app := do.MustInvoke[*codis.HTTPAppService](injector)

	go app.ListenAndServe()

	println("Application started")
	_, err := injector.ShutdownOnSignals(syscall.SIGINT, syscall.SIGTERM)
	if err != nil {
		println(err)
	}

	// Graceful shutdown: stop all consumers
	if err := consumerManager.StopAll(); err != nil {
		slog.Error("Error stopping consumers", "error", err)
	}

	// Close RabbitMQ connection
	if err := connection.Close(); err != nil {
		slog.Error("Error closing RabbitMQ connection", "error", err)
	}

	println("Application stopped")
}
