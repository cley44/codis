package handlerAPIWorkflow

import "codis/models"

type WorkflowsUpdateRequest struct {
	StartingNodesIDs      []string                  `json:"starting_nodes_ids"`
	StartingDiscordEvents []models.DiscordEventType `json:"starting_discord_events"`
	Nodes                 []NodeUpdateRequest       `json:"nodes"`
}

type NodeUpdateRequest struct {
	ID         string                 `json:"id"`
	WorkflowID string                 `json:"workflow_id"`
	Type       models.DiscordNodeType `json:"type"`
	NextNodeID *string                `json:"next_node_id"`
	Data       models.NodeData        `json:"data"`
}
