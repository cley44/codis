package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	"codis/domain/auth"
	authHandlers "codis/handlers/auth"
	controllerDiscord "codis/handlers/discord"
	"codis/handlers/middleware"
	handlerWorkflow "codis/handlers/workflow"
	"codis/utils"
)

type APIRouterService struct {
	discordAPIController  *controllerDiscord.DiscordAPIController
	authAPIController     *authHandlers.AuthAPIController
	sessionService        *auth.SessionService
	workflowAPIController *handlerWorkflow.WorkflowsAPIController
}

func NewAPIRouterService(injector do.Injector) (*APIRouterService, error) {
	m := APIRouterService{
		discordAPIController:  do.MustInvoke[*controllerDiscord.DiscordAPIController](injector),
		sessionService:        do.MustInvoke[*auth.SessionService](injector),
		authAPIController:     do.MustInvoke[*authHandlers.AuthAPIController](injector),
		workflowAPIController: do.MustInvoke[*handlerWorkflow.WorkflowsAPIController](injector),
	}

	m.init()

	return &m, nil
}

func (svc *APIRouterService) init() {

}

func (svc *APIRouterService) RegisterDiscordRoutes(router *gin.Engine, authRouter *gin.RouterGroup) {
	discordAPI := router.Group("/discord")
	{
		discordAPI.GET("/invite_link", svc.discordAPIController.HandleDiscordInviteLink)
		discordAPI.POST("/callback", middleware.ValidateBodyMiddleware(controllerDiscord.DiscordCallbackRequest{}), svc.discordAPIController.HandleDiscordCallback)
	}
	discordAPIAuth := authRouter.Group("/discord")
	{
		discordAPIAuth.GET("/guilds", svc.discordAPIController.HandleDiscordGetGuilds)
		authRouter.GET("/profil", svc.authAPIController.GetProfile)
	}

	// Workflow routes (require auth)
	workflowController := svc.workflowAPIController
	workflows := authRouter.Group("/workflows")
	{
		workflows.GET("", workflowController.ListWorkflows)
		workflows.POST("", middleware.ValidateBodyMiddleware(handlerWorkflow.WorkflowsCreateRequest{}), workflowController.CreateWorkflow)
		workflows.GET(":workflow_id", workflowController.GetWorkflow)
		workflows.PUT(":workflow_id", middleware.ValidateBodyMiddleware(handlerWorkflow.WorkflowsUpdateRequest{}), workflowController.UpdateWorkflow)
		workflows.DELETE(":workflow_id", workflowController.DeleteWorkflow)

		workflows.GET(":workflow_id/nodes", workflowController.ListNodes)
		workflows.POST(":workflow_id/nodes", workflowController.CreateNode)
		workflows.PUT(":workflow_id/nodes/:node_id", workflowController.UpdateNode)
		workflows.DELETE(":workflow_id/nodes/:node_id", workflowController.DeleteNode)
	}

}

func (svc *APIRouterService) RegisterRoutes(router *gin.Engine) {
	// Logout endpoint - accessible without auth middleware
	router.POST("/logout", svc.authAPIController.Logout)
	router.GET("/helloworld", func(ctx *gin.Context) {

		session := sessions.Default(ctx)

		res := session.Get("user_id")

		println(res)
		utils.PrintJSONIndent(res)

		ctx.JSON(200, res)
	})
}
