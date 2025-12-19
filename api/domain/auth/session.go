package auth

import (
	"codis/config"
	"codis/models"
	"codis/repository"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type SessionService struct {
	config                  *config.ConfigService
	userRepository          *repository.UserRepository
	postgresDatabaseService *repository.PostgresDatabaseService
	Store                   postgres.Store
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

	return sessions.Sessions("mysession", svc.Store)
}

func (svc SessionService) GetCurrentUser(userID string) (models.User, error) {

	return svc.userRepository.GetByID(userID)
}
