package handlerAPIAuth

import (
	"codis/domain/auth"
	"codis/utils"
	"errors"
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
		utils.AbortRequest(ctx, http.StatusUnauthorized, errors.New("Unauthorized"), "Unauthorized")
		return
	}

	user, err := svc.sessionService.GetCurrentUser(userID.(string))
	if err != nil {
		utils.AbortRequest(ctx, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	ctx.JSON(200, user)
}
