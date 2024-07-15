package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	controllerDiscord "codis/handlers/discord"
)

type APIRouterService struct {
	discordAPIControllersService *controllerDiscord.DiscordAPIControllersService
}

func NewAPIRouterService(injector do.Injector) (*APIRouterService, error) {
	m := APIRouterService{
		discordAPIControllersService: do.MustInvoke[*controllerDiscord.DiscordAPIControllersService](injector),
	}

	m.init()

	return &m, nil
}

func (svc *APIRouterService) init() {

}

func (svc *APIRouterService) RegisterDiscordRoutes(router *gin.Engine) {
	discordAPI := router.Group("/discord")
	{
		discordAPI.GET("/invite_link", svc.discordAPIControllersService.HandleDiscordInviteLink)
		discordAPI.POST("/callback", svc.discordAPIControllersService.HandleDiscordCallback)
	}
}

func (svc *APIRouterService) RegisterRoutes(router *gin.Engine) {
	router.GET("/helloworld", func(ctx *gin.Context) {
		println("test")
	})
}
