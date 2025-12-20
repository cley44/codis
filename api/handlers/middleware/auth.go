package middleware

import (
	"codis/domain/auth"
	"codis/utils"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AuthMiddlewareService struct {
	sessionService *auth.SessionService
}

func NewAuthMiddlewareService(injector do.Injector) (*AuthMiddlewareService, error) {
	a := AuthMiddlewareService{
		sessionService: do.MustInvoke[*auth.SessionService](injector),
	}

	return &a, nil
}

func (svc AuthMiddlewareService) AuthSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userID := session.Get("user_id")
		if userID == nil {
			utils.AbortRequest(c, http.StatusUnauthorized, errors.New("Unauthorized"), "Unauthorized")
			return
		}

		user, err := svc.sessionService.GetCurrentUser(userID.(string))
		if err != nil {
			utils.AbortRequest(c, http.StatusUnauthorized, err, "Unauthorized")
			return
		}

		c.Set("user_id", userID)
		c.Set("user", user)
		c.Next()
	}
}
