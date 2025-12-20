package handlerAPIDiscord

import (
	"codis/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svc *DiscordAPIController) HandleDiscordGetGuilds(ctx *gin.Context) {

	user, err := svc.sessionService.GetCurrentUserFromContext(ctx)
	if err == nil {
		utils.AbortRequest(ctx, http.StatusUnauthorized, errors.New("Unauthorized"), "Unauthorized")
		return
	}

	guilds, err := svc.discordService.GetGuildsList(user.ID.String(), *user.DiscordSession)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to get guilds list")
		return
	}
	utils.PrintJSONIndent(guilds)
	ctx.JSON(200, guilds)
}
