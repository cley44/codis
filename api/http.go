package codis

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	"codis/domain/auth"
	apiHandlers "codis/handlers"
	middleware "codis/handlers/middleware"
)

type HTTPAppService struct {
	injector         do.Injector
	apiRouterService *apiHandlers.APIRouterService

	router *gin.Engine
}

func NewHTTPAppService(injector do.Injector) (*HTTPAppService, error) {
	s := HTTPAppService{
		injector: injector,

		router: nil,
	}

	return &s, nil
}

func (svc *HTTPAppService) ShutDown() error {
	return nil
}

func (svc *HTTPAppService) ListenAndServe() {
	svc.router = gin.New()
	// @TODO change this to allow only the correct host via config
	svc.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true, // 🔴 REQUIRED
	}))

	// Recover middleware
	svc.router.Use(middleware.GinOopsRecovery())

	// Init sessions middleware
	s := do.MustInvoke[*auth.SessionService](svc.injector)
	svc.router.Use(s.InitSessionMiddleware())

	svc.apiRouterService = do.MustInvoke[*apiHandlers.APIRouterService](svc.injector)

	authRouter := svc.router.Group("")
	authRouter.Use(middleware.AuthSessionMiddleware())

	svc.apiRouterService.RegisterRoutes(svc.router)
	svc.apiRouterService.RegisterDiscordRoutes(svc.router, authRouter)

	err := svc.router.Run(":80")
	if err != nil {
		println(err)
	}
}
