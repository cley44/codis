package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	"codis/domain/auth"
	authHandlers "codis/handlers/auth"
	controllerDiscord "codis/handlers/discord"
	"codis/utils"
)

type APIRouterService struct {
	discordAPIController *controllerDiscord.DiscordAPIController
	authAPIController    *authHandlers.AuthAPIController
	sessionService       *auth.SessionService
}

func NewAPIRouterService(injector do.Injector) (*APIRouterService, error) {
	m := APIRouterService{
		discordAPIController: do.MustInvoke[*controllerDiscord.DiscordAPIController](injector),
		sessionService:       do.MustInvoke[*auth.SessionService](injector),
		authAPIController:    do.MustInvoke[*authHandlers.AuthAPIController](injector),
	}

	m.init()

	return &m, nil
}

func (svc *APIRouterService) init() {

}

func (svc *APIRouterService) RegisterDiscordRoutes(router *gin.Engine, authRouter *gin.RouterGroup) {
	discordAPI := router.Group("/discord")
	{
		discordAPI.GET("/invite_link", svc.discordAPIController.HandleDiscordInviteLink)
		discordAPI.POST("/callback", svc.discordAPIController.HandleDiscordCallback)
	}
	discordAPIAuth := authRouter.Group("/discord")
	{
		discordAPIAuth.GET("/guilds", svc.discordAPIController.HandleDiscordGetGuilds)
		authRouter.GET("/profil", svc.authAPIController.GetProfile)
	}
}

func (svc *APIRouterService) RegisterRoutes(router *gin.Engine) {
	router.GET("/helloworld", func(ctx *gin.Context) {

		session := sessions.Default(ctx)

		res := session.Get("user_id")

		println(res)
		utils.PrintJSONIndent(res)

		ctx.JSON(200, res)
	})
}
