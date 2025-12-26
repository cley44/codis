package discordHandler

import (
	"codis/domain/rabbitmq"
	"codis/models"

	"github.com/samber/do/v2"
)

type NodeHandlerService struct {
	HandlerMap map[models.DiscordNodeType]rabbitmq.NodeHandler
}

func NewNodeHandlerService(injector do.Injector) (*NodeHandlerService, error) {
	handlerMap := make(map[models.DiscordNodeType]rabbitmq.NodeHandler)

	handlers := []rabbitmq.NodeHandler{
		do.MustInvoke[*HandlerAddMemberRole](injector),
		// Add more handlers here
	}

	for _, handler := range handlers {
		handlerMap[handler.GetType()] = handler
	}

	return &NodeHandlerService{
		HandlerMap: handlerMap,
	}, nil
}

func (h *NodeHandlerService) GetHandler(actionType models.DiscordNodeType) (rabbitmq.NodeHandler, bool) {
	handler, exists := h.HandlerMap[actionType]
	return handler, exists
}
