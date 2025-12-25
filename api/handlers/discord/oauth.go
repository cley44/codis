package handlerAPIDiscord

import (
	"codis/handlers/handlers"
	"codis/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (svc *DiscordAPIController) HandleDiscordInviteLink(ctx *gin.Context) {
	inviteLink := svc.discordService.GetDiscordInviteLink()

	ctx.JSON(200, gin.H{
		"discord_invite_link": inviteLink,
	})
}

func (svc *DiscordAPIController) HandleDiscordCallback(ctx *gin.Context) {
	body := handlers.GetBody(ctx).(*DiscordCallbackRequest)

	discordOauthSession := svc.discordService.StartSession(body.Code, body.State)

	discordUser, err := svc.discordService.GetUser(discordOauthSession.Session)
	if err != nil {
		utils.AbortRequest(ctx, 500, err, "Failed to handle discord callback (0)")
		return
	}

	user, err := svc.userRepository.CreateOrUpdate(
		discordUser.Username,
		discordUser.GlobalName,
		discordUser.ID.String(),
		discordUser.Avatar,
		&discordOauthSession,
		discordUser.Email)
	if err != nil {
		utils.AbortRequest(ctx, 500, err, "Failed to handle discord callback (1)")
		return
	}

	session := sessions.Default(ctx)

	session.Set("user_id", user.ID.String())
	err = session.Save()
	if err != nil {
		utils.AbortRequest(ctx, 500, err, "Failed to handle discord callback (2)")
		return
	}

	ctx.JSON(200, user)
}
