package workers

import (
	"codis/domain/rabbitmq"
	"codis/models"
	"codis/repository"
	"log/slog"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type DispatchWorker struct {
	// Add your dependencies here (e.g., repositories, services, etc.)
	workflowRepository *repository.WorkflowRepository
	publisherService   *rabbitmq.PublisherService
}

func NewDispatchWorker(injector do.Injector) *DispatchWorker {
	return &DispatchWorker{
		workflowRepository: do.MustInvoke[*repository.WorkflowRepository](injector),
		publisherService:   do.MustInvoke[*rabbitmq.PublisherService](injector),
	}
}

// QueueName returns the name of the queue this worker consumes from
func (w *DispatchWorker) QueueName() rabbitmq.RoutingKey {
	return rabbitmq.RoutingKeyDispatch
}

// QueueOptions returns the queue configuration options
func (w *DispatchWorker) QueueOptions() *rabbitmq.QueueOptions {
	return &rabbitmq.QueueOptions{
		Durable:    true,  // Queue survives broker restart
		AutoDelete: false, // Queue is not deleted when no longer used
		Exclusive:  false, // Queue can be accessed by other connections
		NoWait:     false, // Wait for server confirmation
		Args:       nil,   // Additional queue arguments
	}
}

// HandleMessage processes a single message from the queue
func (w *DispatchWorker) HandleMessage(msg rabbitmq.AMQPMessage) error {
	slog.Info("Processing message",
		"queue", w.QueueName(),
		"routing_key", msg.RoutingKey,
	)

	// Example: Process the message
	slog.Info("Message payload", "payload", msg.Body)
	// utils.PrintJSONIndent(msg.Body.DiscordEvent.MessageCreateEvent.Message.Member.User.Username)

	workflow, err := w.workflowRepository.GetByStartingDiscordEvents(msg.Body.DiscordEvent.GuildID, []models.DiscordEventType{msg.Body.DiscordEvent.Type})
	if err != nil {
		return oops.Wrapf(err, "Failed to get workflow")
	}
	if workflow.StartingNodesIDs == nil {
		return nil
	}

	for _, nodeID := range *workflow.StartingNodesIDs {
		msg.Body.DiscordEvent.NodeIDToExecute = nodeID
		err = w.publisherService.Publish(rabbitmq.RoutingKeyNodeExecute, msg.Body)
		if err != nil {
			return oops.Wrapf(err, "Failed to publish message to example queue")
		}
	}

	// Example: Do some work
	// result := w.processData(payload)

	// If processing fails, return an error
	// The message will be nacked and requeued
	// If processing succeeds, return nil
	// The message will be acked

	return nil
}
