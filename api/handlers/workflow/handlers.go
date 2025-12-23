package handlerAPIWorkflow

import (
	"encoding/json"
	"codis/models"
	"codis/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// ListWorkflows handles GET /workflows
func (svc *WorkflowsAPIController) ListWorkflows(ctx *gin.Context) {
	guildID := ctx.Query("guild_id")
	var gid *string
	if guildID != "" {
		gid = &guildID
	}

	workflows, err := svc.workflowRepository.List(gid)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to list workflows")
		return
	}

	ctx.JSON(200, workflows)
}

// CreateWorkflow handles POST /workflows
func (svc *WorkflowsAPIController) CreateWorkflow(ctx *gin.Context) {
	var body struct {
		StartNodeID           *string               `json:"start_node_id"`
		GuildID               string                `json:"guild_id" binding:"required"`
		StartingDiscordEvents []models.DiscordEventType `json:"starting_discord_events"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid body")
		return
	}

	wf, err := svc.workflowRepository.Create(
		lo.Ternary(body.StartNodeID != nil, *body.StartNodeID, ""),
		body.GuildID,
		body.StartingDiscordEvents,
	)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to create workflow")
		return
	}

	ctx.JSON(http.StatusCreated, wf)
}

// GetWorkflow handles GET /workflows/:id
func (svc *WorkflowsAPIController) GetWorkflow(ctx *gin.Context) {
	id := ctx.Param("workflow_id")
	wf, err := svc.workflowRepository.GetByID(id)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusNotFound, err, "Workflow not found")
		return
	}

	ctx.JSON(200, wf)
}

// UpdateWorkflow handles PUT /workflows/:id
func (svc *WorkflowsAPIController) UpdateWorkflow(ctx *gin.Context) {
	id := ctx.Param("workflow_id")
	var body struct {
		StartNodeID           *string               `json:"start_node_id"`
		StartingDiscordEvents []models.DiscordEventType `json:"starting_discord_events"`
		Edges                 map[string]interface{} `json:"edges"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid body")
		return
	}

	wf, err := svc.workflowRepository.GetByID(id)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusNotFound, err, "Workflow not found")
		return
	}

	if body.StartNodeID != nil {
		wf.StartNodeID = *body.StartNodeID
	}

	if body.StartingDiscordEvents != nil {
		wf.StartingDiscordEvents = body.StartingDiscordEvents
	}

	// marshal edges back into json.RawMessage
	if body.Edges != nil {
		b, err := json.Marshal(body.Edges)
		if err != nil {
			utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid edges")
			return
		}
		wf.Edges = b
	}

	updated, err := svc.workflowRepository.Update(wf)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to update workflow")
		return
	}

	ctx.JSON(200, updated)
}

// DeleteWorkflow handles DELETE /workflows/:id
func (svc *WorkflowsAPIController) DeleteWorkflow(ctx *gin.Context) {
	id := ctx.Param("workflow_id")
	err := svc.workflowRepository.Delete(id)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to delete workflow")
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ListNodes handles GET /workflows/:workflow_id/nodes
func (svc *WorkflowsAPIController) ListNodes(ctx *gin.Context) {
	workflowID := ctx.Param("workflow_id")
	nodes, err := svc.nodeRepository.ListByWorkflow(workflowID)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to list nodes")
		return
	}

	ctx.JSON(200, nodes)
}

// CreateNode handles POST /workflows/:workflow_id/nodes
func (svc *WorkflowsAPIController) CreateNode(ctx *gin.Context) {
	workflowID := ctx.Param("workflow_id")
	var body struct {
		Type       string                 `json:"type" binding:"required"`
		NextNodeID *string                `json:"next_node_id"`
		Meta       map[string]interface{} `json:"meta"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid body")
		return
	}

	metaRaw := []byte("{}")
	if body.Meta != nil {
		b, err := json.Marshal(body.Meta)
		if err != nil {
			utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid meta")
			return
		}
		metaRaw = b
	}

	node := models.Node{
		WorkflowID: workflowID,
		Type:       body.Type,
		NextNodeID: body.NextNodeID,
		Meta:       metaRaw,
	}

	created, err := svc.nodeRepository.Create(node)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to create node")
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

// UpdateNode handles PUT /workflows/:workflow_id/nodes/:node_id
func (svc *WorkflowsAPIController) UpdateNode(ctx *gin.Context) {
	var body struct {
		Type       *string                `json:"type"`
		NextNodeID *string                `json:"next_node_id"`
		Meta       map[string]interface{} `json:"meta"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid body")
		return
	}

	nodeID := ctx.Param("node_id")
	node, err := svc.nodeRepository.GetByID(nodeID)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusNotFound, err, "Node not found")
		return
	}

	if body.Type != nil {
		node.Type = *body.Type
	}

	if body.NextNodeID != nil {
		node.NextNodeID = body.NextNodeID
	}

	if body.Meta != nil {
		b, err := json.Marshal(body.Meta)
		if err != nil {
			utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid meta")
			return
		}
		node.Meta = b
	}

	updated, err := svc.nodeRepository.Update(node)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to update node")
		return
	}

	ctx.JSON(200, updated)
}

// DeleteNode handles DELETE /workflows/:workflow_id/nodes/:node_id
func (svc *WorkflowsAPIController) DeleteNode(ctx *gin.Context) {
	nodeID := ctx.Param("node_id")
	err := svc.nodeRepository.Delete(nodeID)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to delete node")
		return
	}

	ctx.Status(http.StatusNoContent)
}
