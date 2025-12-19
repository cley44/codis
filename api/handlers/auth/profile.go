package handlerAPIAuth

import (
	"codis/domain/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type AuthAPIController struct {
	sessionService *auth.SessionService
}

func NewAuthAPIController(injector do.Injector) (*AuthAPIController, error) {
	auth := AuthAPIController{
		sessionService: do.MustInvoke[*auth.SessionService](injector),
	}

	return &auth, nil
}

func (svc *AuthAPIController) GetProfile(ctx *gin.Context) {

	userID, exist := ctx.Get("user_id")
	if !exist {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "authentificaton required"},
		)
	}

	user, err := svc.sessionService.GetCurrentUser(userID.(string))
	if err != nil {
		panic(err)
	}

	ctx.JSON(200, user)
}
