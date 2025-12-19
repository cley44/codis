package handlerAPIDiscord

import (
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
	var body DiscordCallbackRequest
	if err := ctx.BindJSON(&body); err != nil {
		panic(err)
	}

	discordOauthSession := svc.discordService.StartSession(body.Code, body.State)

	discordUser, err := svc.discordService.GetUser(discordOauthSession.Session)
	if err != nil {
		utils.PrintJSONIndent(err.Error())
		//@TODO should be handled
		panic(err.Error())
	}

	user, err := svc.userRepository.Create(
		discordUser.Username,
		discordUser.GlobalName,
		discordUser.ID.String(),
		discordUser.Avatar,
		&discordOauthSession,
		discordUser.Email)
	if err != nil {
		utils.PrintJSONIndent(err.Error())
		ctx.JSON(200, "failed")
		return
	}

	session := sessions.Default(ctx)

	session.Set("user_id", user.ID.String())
	err = session.Save()
	if err != nil {
		panic(err)
	}

	ctx.JSON(200, user)
}
