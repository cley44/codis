package discordHandler

import (
	discordHandlerMessage "codis/domain/discord/handlers/message"
	discordHandlerRole "codis/domain/discord/handlers/role"
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
		do.MustInvoke[*discordHandlerRole.HandlerAddMemberRole](injector),
		do.MustInvoke[*discordHandlerRole.HandlerRemoveMemberRole](injector),
		do.MustInvoke[*discordHandlerMessage.HandlerSendMessage](injector),
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
