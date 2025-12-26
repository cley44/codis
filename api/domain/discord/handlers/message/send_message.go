package discordHandlerMessage

import (
	"codis/domain/discord"
	"codis/domain/rabbitmq"
	"codis/models"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type HandlerSendMessage struct {
	discordService *discord.DiscordService
}

func NewHandlerSendMessage(injector do.Injector) (*HandlerSendMessage, error) {
	return &HandlerSendMessage{
		discordService: do.MustInvoke[*discord.DiscordService](injector),
	}, nil
}

func (h *HandlerSendMessage) GetType() models.DiscordNodeType {
	return models.DiscordNodeTypeSendMessage
}

func (h *HandlerSendMessage) Execute(msg rabbitmq.AMQPMessageBody, node models.Node) error {

	if node.Data.ChannelID == nil {
		return oops.With("node_id", node.ID).Errorf("Channel ID is null in node data")
	}
	if node.Data.MessageContent == nil {
		return oops.With("node_id", node.ID).Errorf("Message content is null in node data")
	}
	channelID := *node.Data.ChannelID
	messageContent := *node.Data.MessageContent

	err := h.discordService.SendMessage(channelID, messageContent)

	return err
}
