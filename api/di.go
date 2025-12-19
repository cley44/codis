package codis

import (
	"codis/config"
	"codis/domain/auth"
	"codis/domain/discord"
	"codis/handlers"
	handlerAPIAuth "codis/handlers/auth"
	handlerAPIDiscord "codis/handlers/discord"
	"codis/instrumentation"
	"codis/repository"

	"github.com/samber/do/v2"
)

func RegisterDatabase(injector do.Injector) {
	do.Provide(injector, repository.NewPostgresDatabaseService)

	RegisterRepository(injector)
}

func RegisterRepository(injector do.Injector) {
	do.Provide(injector, repository.NewUserRepository)
}

func RegisterControllers(injector do.Injector) {
	do.Provide(injector, handlerAPIDiscord.NewDiscordAPIControllersService)
	do.Provide(injector, handlerAPIAuth.NewAuthAPIController)
}

func RegisterBase(injector do.Injector) {
	do.Provide(injector, config.NewConfigService)
}

func RegisterDiscord(injector do.Injector) {
	do.Provide(injector, discord.NewDiscordService)
}

func RegisterAPI(injector do.Injector) {
	do.Provide(injector, handlers.NewAPIRouterService)

	do.Provide(injector, auth.NewSessionService)

	RegisterControllers(injector)
}

func RegisterApp(injector do.Injector) {
	do.Provide(injector, NewHTTPAppService)
}

func RegisterInstrumentation(injector do.Injector) {
	do.Provide(injector, instrumentation.NewLoggerService)
}

func RegisterAll() *do.RootScope {
	injector := do.New()

	RegisterBase(injector)

	RegisterInstrumentation(injector)

	RegisterAPI(injector)

	RegisterDatabase(injector)

	RegisterApp(injector)
	RegisterDiscord(injector)

	return injector
}
