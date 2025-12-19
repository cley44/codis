package handlerAPIDiscord

import (
	"github.com/gin-gonic/gin"
)

func (svc *DiscordAPIControllersService) HandleDiscordInviteLink(ctx *gin.Context) {
	inviteLink := svc.discordService.GetDiscordInviteLink()

	ctx.JSON(200, gin.H{
		"discord_invite_link": inviteLink,
	})
}

func (svc *DiscordAPIControllersService) HandleDiscordCallback(ctx *gin.Context) {
	var body DiscordCallbackRequest
	if err := ctx.BindJSON(&body); err != nil {
		panic(err)
	}

	svc.discordService.StartSession(body.Code, body.State)

	svc.userRepository.Create("test", "emial", "password")

	// ctx.JSON(200, ctx.json)
}
