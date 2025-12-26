package workers

import (
	discord "codis/domain/discord/handlers"
	"codis/domain/rabbitmq"
	"codis/repository"
	"log/slog"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type NodeExecuteWorker struct {
	// Add your dependencies here (e.g., repositories, services, etc.)
	nodeRepository     *repository.NodeRepository
	publisherService   *rabbitmq.PublisherService
	nodeHandlerService *discord.NodeHandlerService
}

func NewNodeExecuteWorker(injector do.Injector) *NodeExecuteWorker {
	return &NodeExecuteWorker{
		nodeRepository:     do.MustInvoke[*repository.NodeRepository](injector),
		publisherService:   do.MustInvoke[*rabbitmq.PublisherService](injector),
		nodeHandlerService: do.MustInvoke[*discord.NodeHandlerService](injector),
	}
}

// QueueName returns the name of the queue this worker consumes from
func (w *NodeExecuteWorker) QueueName() rabbitmq.RoutingKey {
	return rabbitmq.RoutingKeyNodeExecute
}

// QueueOptions returns the queue configuration options
func (w *NodeExecuteWorker) QueueOptions() *rabbitmq.QueueOptions {
	return &rabbitmq.QueueOptions{
		Durable:    true,  // Queue survives broker restart
		AutoDelete: false, // Queue is not deleted when no longer used
		Exclusive:  false, // Queue can be accessed by other connections
		NoWait:     false, // Wait for server confirmation
		Args:       nil,   // Additional queue arguments
	}
}

// HandleMessage processes a single message from the queue
func (w *NodeExecuteWorker) HandleMessage(msg rabbitmq.AMQPMessage) error {
	slog.Info("Processing message",
		"queue", w.QueueName(),
		"routing_key", msg.RoutingKey,
	)

	// Example: Process the message
	slog.Info("Message payload", "payload", msg.Body)

	node, err := w.nodeRepository.GetByID(msg.Body.DiscordEvent.NodeIDToExecute)
	if err != nil {
		return oops.Wrapf(err, "Failed to get node")
	}

	handler, exists := w.nodeHandlerService.GetHandler(node.Type)
	if !exists {
		return oops.Errorf("No handler found for node type %s", node.Type)
	}

	err = handler.Execute(msg.Body, node)
	if err != nil {
		return oops.With("node_id", node.ID).Wrapf(err, "Failed to execute node handler")
	}

	if node.NextNodeID != nil {
		msg.Body.DiscordEvent.NodeIDToExecute = *node.NextNodeID
		err = w.publisherService.Publish(rabbitmq.RoutingKeyNodeExecute, msg.Body)
		if err != nil {
			return oops.Wrapf(err, "Failed to publish message to node_execute queue")
		}
	}

	// Execute node logic here

	// err = w.publisherService.Publish(rabbitmq.RoutingKeyNodeExecute, msg.Body)
	// if err != nil {
	// 	return oops.Wrapf(err, "Failed to publish message to example queue")
	// }

	// Example: Do some work
	// result := w.processData(payload)

	// If processing fails, return an error
	// The message will be nacked and requeued
	// If processing succeeds, return nil
	// The message will be acked

	return nil
}
