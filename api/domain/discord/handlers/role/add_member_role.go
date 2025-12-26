package discordHandlerRole

import (
	"codis/domain/discord"
	"codis/domain/rabbitmq"
	"codis/models"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type HandlerAddMemberRole struct {
	discordService *discord.DiscordService
}

func NewHandlerAddMemberRole(injector do.Injector) (*HandlerAddMemberRole, error) {
	return &HandlerAddMemberRole{
		discordService: do.MustInvoke[*discord.DiscordService](injector),
	}, nil
}

func (h *HandlerAddMemberRole) GetType() models.DiscordNodeType {
	return models.DiscordNodeTypeAddMemberRole
}

func (h *HandlerAddMemberRole) Execute(msg rabbitmq.AMQPMessageBody, node models.Node) error {

	guildID := msg.DiscordEvent.GuildID
	userID := msg.DiscordEvent.UserID
	if node.Data.RoleID == nil {
		return oops.With("guild_id", guildID).With("user_id", userID).Errorf("Role ID is null in node data")
	}
	roleID := *node.Data.RoleID
	// roleID := "1453810943912972521"

	err := h.discordService.AddMemberRole(guildID, userID, roleID)

	return err
}
