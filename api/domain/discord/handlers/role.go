package discord

import (
	"codis/domain/discord"
	"codis/domain/rabbitmq"
	"codis/models"
)

type HandlerAddMemberRole struct {
	discordService discord.DiscordService
}

func NewHandlerAddMemberRole(discordService discord.DiscordService) *HandlerAddMemberRole {
	return &HandlerAddMemberRole{
		discordService: discordService,
	}
}

func (h *HandlerAddMemberRole) GetType() models.DiscordNodeType {
	return models.DiscordNodeTypeAddMemberRole
}

func (h *HandlerAddMemberRole) Execute(msg rabbitmq.AMQPMessageBody) error {

	guildID := msg.DiscordEvent.GuildID
	userID := msg.DiscordEvent.UserID
	roleID := msg.DiscordEvent.RoleID

	err := h.discordService.AddMemberRole(guildID, userID, roleID)

	return err
}
