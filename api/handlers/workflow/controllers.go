package handlerAPIWorkflow

import (
	"codis/repository"

	"codis/domain/auth"

	"github.com/samber/do/v2"
)

type WorkflowsAPIController struct {
	workflowRepository *repository.WorkflowRepository
	nodeRepository     *repository.NodeRepository
	sessionService     *auth.SessionService
}

func NewWorkflowsAPIController(injector do.Injector) (*WorkflowsAPIController, error) {
	w := WorkflowsAPIController{
		workflowRepository: do.MustInvoke[*repository.WorkflowRepository](injector),
		nodeRepository:     do.MustInvoke[*repository.NodeRepository](injector),
		sessionService:     do.MustInvoke[*auth.SessionService](injector),
	}

	w.init()
	return &w, nil
}

func (svc *WorkflowsAPIController) init() {
	// noop for now
}

// Http handlers follow... (see handlers.go)
