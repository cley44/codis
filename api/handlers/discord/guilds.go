package handlerAPIDiscord

import (
	"codis/utils"
	"errors"
	"net/http"

	"github.com/disgoorg/disgo/discord"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
)

func (svc *DiscordAPIController) HandleDiscordGetGuilds(ctx *gin.Context) {
	user, err := svc.sessionService.GetCurrentUserFromContext(ctx)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusUnauthorized, errors.New("Unauthorized"), "Unauthorized")
		return
	}

	guilds, err := svc.discordService.GetGuildsList(user.ID.String(), *user.DiscordSession)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, oops.Wrap(err), "Failed to get guilds list")
		return
	}

	guildsResponse := []DiscordGuilds{}

	for _, g := range guilds {
		if g.Permissions.Has(discord.PermissionAdministrator) {
			_, exist := svc.discordService.IsBotAMemberOfguild(g.ID)
			guildResponse := DiscordGuilds{
				ID:            g.ID.String(),
				Name:          g.Name,
				BannerURL:     g.Banner,
				IconURL:       g.Icon,
				Owner:         g.Owner,
				BotInviteLink: svc.discordService.GetDiscordGuildInviteLink(g),
				BotPresent:    exist,
			}
			guildsResponse = append(guildsResponse, guildResponse)
		}
	}

	ctx.JSON(200, guildsResponse)
}
