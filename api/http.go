package codis

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	apiHandlers "codis/handlers"
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

	svc.apiRouterService = do.MustInvoke[*apiHandlers.APIRouterService](svc.injector)

	svc.apiRouterService.RegisterRoutes(svc.router)
	svc.apiRouterService.RegisterDiscordRoutes(svc.router)

	err := svc.router.Run(":80")
	if err != nil {
		println(err)
	}
}
