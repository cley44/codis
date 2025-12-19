package handlerAPIDiscord

import (
	"codis/domain/auth"
	"codis/domain/discord"
	"codis/repository"

	"github.com/samber/do/v2"
)

type DiscordAPIController struct {
	discordService *discord.DiscordService
	userRepository *repository.UserRepository
	sessionService *auth.SessionService
}

func NewDiscordAPIControllersService(injector do.Injector) (*DiscordAPIController, error) {
	m := DiscordAPIController{
		discordService: do.MustInvoke[*discord.DiscordService](injector),
		userRepository: do.MustInvoke[*repository.UserRepository](injector),
		sessionService: do.MustInvoke[*auth.SessionService](injector),
	}

	m.init()

	return &m, nil
}

func (svc *DiscordAPIController) init() {

}
