package handlerAPIDiscord

import (
	"codis/domain/discord"
	"codis/repository"

	"github.com/samber/do/v2"
)

type DiscordAPIControllersService struct {
	discordService *discord.DiscordService
	userRepository *repository.UserRepository
}

func NewDiscordAPIControllersService(injector do.Injector) (*DiscordAPIControllersService, error) {
	m := DiscordAPIControllersService{
		discordService: do.MustInvoke[*discord.DiscordService](injector),
		userRepository: do.MustInvoke[*repository.UserRepository](injector),
	}

	m.init()

	return &m, nil
}

func (svc *DiscordAPIControllersService) init() {

}
