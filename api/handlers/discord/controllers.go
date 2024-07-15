package handlerAPIDiscord

import (
	"codis/discord"

	"github.com/samber/do/v2"
)

type DiscordAPIControllersService struct {
	discordService *discord.DiscordService
}

func NewDiscordAPIControllersService(injector do.Injector) (*DiscordAPIControllersService, error) {
	m := DiscordAPIControllersService{
		discordService: do.MustInvoke[*discord.DiscordService](injector),
	}

	m.init()

	return &m, nil
}

func (svc *DiscordAPIControllersService) init() {

}
