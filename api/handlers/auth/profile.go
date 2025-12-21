package handlerAPIAuth

import (
	"codis/domain/auth"
	"codis/utils"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
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
	user, err := svc.sessionService.GetCurrentUserFromContext(ctx)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusUnauthorized, errors.New("Unauthorized"), "Unauthorized")
		return
	}

	ctx.JSON(200, user)
}

func (svc *AuthAPIController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	
	// Clear the session
	session.Clear()
	err := session.Save()
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to logout")
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Logged out successfully",
	})
}
