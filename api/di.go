package codis

import (
	"codis/config"
	"codis/discord"
	"codis/handlers"
	handlerAPIDiscord "codis/handlers/discord"
	"codis/repository"

	"github.com/samber/do/v2"
)

func RegisterDatabase(injector do.Injector) {
	do.Provide(injector, repository.NewPostgresDatabaseService)
}

func RegisterControllers(injector do.Injector) {
	do.Provide(injector, handlerAPIDiscord.NewDiscordAPIControllersService)
}

func RegisterBase(injector do.Injector) {
	do.Provide(injector, config.NewConfigService)
}

func RegisterDiscord(injector do.Injector) {
	do.Provide(injector, discord.NewDiscordService)
}

func RegisterAPI(injector do.Injector) {
	do.Provide(injector, handlers.NewAPIRouterService)

	RegisterControllers(injector)
}

func RegisterApp(injector do.Injector) {
	do.Provide(injector, NewHTTPAppService)
}

func RegisterAll() *do.RootScope {
	injector := do.New()

	RegisterBase(injector)

	RegisterAPI(injector)

	RegisterDatabase(injector)

	RegisterApp(injector)
	RegisterDiscord(injector)

	return injector
}
