package auth

import (
	"codis/config"
	"codis/domain/discord"
	"codis/models"
	"codis/repository"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type SessionService struct {
	config                  *config.ConfigService
	userRepository          *repository.UserRepository
	postgresDatabaseService *repository.PostgresDatabaseService
	Store                   postgres.Store
	discordService          *discord.DiscordService
}

func NewSessionService(injector do.Injector) (*SessionService, error) {
	config := do.MustInvoke[*config.ConfigService](injector)
	postgresDatabaseService := do.MustInvoke[*repository.PostgresDatabaseService](injector)

	store, err := postgres.NewStore(postgresDatabaseService.Db.DB, []byte(config.Auth.SessionSecret))
	if err != nil {
		panic(err)
	}

	s := SessionService{
		config:                  config,
		userRepository:          do.MustInvoke[*repository.UserRepository](injector),
		postgresDatabaseService: postgresDatabaseService,
		Store:                   store,
		discordService:          do.MustInvoke[*discord.DiscordService](injector),
	}

	return &s, nil
}

func (svc SessionService) InitSessionMiddleware() gin.HandlerFunc {

	svc.Store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
		Secure:   false, // true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,
	})

	return sessions.Sessions("codis_session", svc.Store)
}

func (svc SessionService) GetCurrentUser(userID string) (models.User, error) {
	return svc.userRepository.GetByID(userID)
}

func (svc SessionService) GetCurrentUserFromContext(ctx *gin.Context) (models.User, error) {
	userContext, exist := ctx.Get("user")
	if !exist {
		return models.User{}, oops.Errorf("No user in context")
	}

	user, ok := userContext.(models.User)
	if !ok {
		return models.User{}, oops.Errorf("Failed to cast user to models")
	}

	// We verify that the discord oauth session is still valid if not we refresh it
	session, err := svc.discordService.VerifySession(user.ID.String(), *user.DiscordSession)
	if err != nil {
		return models.User{}, oops.Wrap(err)
	}
	user.DiscordSession = &session

	return user, nil
}
