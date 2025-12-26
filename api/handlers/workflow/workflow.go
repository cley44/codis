package handlerAPIWorkflow

import (
	"codis/handlers/handlers"
	"codis/models"
	"codis/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// ListWorkflows handles GET /workflows
func (svc *WorkflowsAPIController) ListWorkflows(ctx *gin.Context) {
	guildID := ctx.Query("guild_id")

	workflows, err := svc.workflowRepository.ListByGuildID(guildID, true, true)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to list workflows")
		return
	}

	ctx.JSON(200, workflows)
}

// CreateWorkflow handles POST /workflows
func (svc *WorkflowsAPIController) CreateWorkflow(ctx *gin.Context) {
	var body struct {
		StartingNodesIDs      []string                  `json:"starting_nodes_ids"`
		GuildID               string                    `json:"guild_id" binding:"required"`
		StartingDiscordEvents []models.DiscordEventType `json:"starting_discord_events"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid body")
		return
	}

	wf, err := svc.workflowRepository.Create(
		lo.Ternary(body.StartingNodesIDs != nil, body.StartingNodesIDs, []string{}),
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
	wf, err := svc.workflowRepository.GetByID(id, true, true)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusNotFound, err, "Workflow not found")
		return
	}

	ctx.JSON(200, wf)
}

// UpdateWorkflow handles PUT /workflows/:id
func (svc *WorkflowsAPIController) UpdateWorkflow(ctx *gin.Context) {
	id := ctx.Param("workflow_id")

	body := handlers.GetBody(ctx).(*WorkflowsUpdateRequest)

	nodes, err := svc.nodeRepository.ListByWorkflowID(id)
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to list nodes")
		return
	}

	existingNodesIDs := lo.Map(nodes, func(n models.Node, _ int) string {
		return n.ID
	})
	incomingNodesIDs := lo.Map(body.Nodes, func(n NodeUpdateRequest, _ int) string {
		return n.ID
	})

	// Determine nodes to delete
	nodesToDelete, _ := lo.Difference(existingNodesIDs, incomingNodesIDs)

	toUpdateNodes := lo.Filter(body.Nodes, func(n NodeUpdateRequest, _ int) bool {
		return lo.Contains(existingNodesIDs, n.ID)
	})

	toInsertNodes := lo.Filter(body.Nodes, func(n NodeUpdateRequest, _ int) bool {
		return !lo.Contains(existingNodesIDs, n.ID)
	})

	// Delete nodes
	if len(nodesToDelete) > 0 {
		err = svc.nodeRepository.Delete(nodesToDelete)
		if err != nil {
			utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to delete nodes")
			return
		}
	}
	if len(toInsertNodes) > 0 {
		// Insert new nodes
		_, err = svc.nodeRepository.CreateMany(lo.Map(toInsertNodes, func(n NodeUpdateRequest, _ int) models.Node {
			return models.Node{
				ID:         n.ID,
				WorkflowID: id,
				Type:       n.Type,
				NextNodeID: n.NextNodeID,
				Data:       n.Data,
			}
		}))
		if err != nil {
			utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to create nodes")
			return
		}
	}

	// Update existing nodes
	for _, n := range toUpdateNodes {
		node := models.Node{
			ID:         n.ID,
			Type:       n.Type,
			NextNodeID: n.NextNodeID,
			Data:       n.Data,
		}
		_, err := svc.nodeRepository.Update(node)
		if err != nil {
			utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to update node")
			return
		}
	}

	updated, err := svc.workflowRepository.Update(id, body.StartingDiscordEvents, body.StartingNodesIDs)
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
	nodes, err := svc.nodeRepository.ListByWorkflowID(workflowID)
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
		Type       models.DiscordNodeType `json:"type" binding:"required"`
		NextNodeID *string                `json:"next_node_id"`
		Data       models.NodeData        `json:"data"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		utils.AbortRequest(ctx, http.StatusBadRequest, err, "Invalid body")
		return
	}

	node := models.Node{
		WorkflowID: workflowID,
		Type:       body.Type,
		NextNodeID: body.NextNodeID,
		Data:       body.Data,
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
		Type       *models.DiscordNodeType `json:"type"`
		NextNodeID *string                 `json:"next_node_id"`
		Data       *models.NodeData        `json:"data"`
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

	if body.Data != nil {
		node.Data = *body.Data
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
	err := svc.nodeRepository.Delete([]string{nodeID})
	if err != nil {
		utils.AbortRequest(ctx, http.StatusInternalServerError, err, "Failed to delete node")
		return
	}

	ctx.Status(http.StatusNoContent)
}
