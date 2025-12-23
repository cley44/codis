package codis

import (
	"codis/config"
	"codis/domain/auth"
	"codis/domain/discord"
	rabbitmqDomain "codis/domain/rabbitmq"
	"codis/handlers"
	handlerAPIAuth "codis/handlers/auth"
	handlerAPIDiscord "codis/handlers/discord"
	"codis/handlers/middleware"
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
	do.Provide(injector, middleware.NewAuthMiddlewareService)
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

func RegisterRabbitMQ(injector do.Injector) {
	// Register in order: Connection → QueueManager → ConsumerManager
	do.Provide(injector, rabbitmqDomain.NewRabbitMQConnectionService)
	do.Provide(injector, rabbitmqDomain.NewQueueManagerService)
	do.Provide(injector, rabbitmqDomain.NewConsumerManagerService)
	do.Provide(injector, rabbitmqDomain.NewPublisherService)
}

func RegisterAll() *do.RootScope {
	injector := do.New()

	RegisterBase(injector)

	RegisterInstrumentation(injector)

	RegisterAPI(injector)

	RegisterDatabase(injector)

	RegisterRabbitMQ(injector)

	RegisterApp(injector)
	RegisterDiscord(injector)

	return injector
}
