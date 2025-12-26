package repository

import (
	"codis/models"

	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type WorkflowRepository struct {
	postgresDatabaseService *PostgresDatabaseService
}

func NewWorkflowRepository(injector do.Injector) (*WorkflowRepository, error) {
	u := WorkflowRepository{
		postgresDatabaseService: do.MustInvoke[*PostgresDatabaseService](injector),
	}

	return &u, nil
}

// Create inserts a new workflow into the database and returns the created workflow.
func (w WorkflowRepository) Create(startingNodesIDs []string, guildID string, startingDiscordEvents []models.DiscordEventType) (workflow models.Workflow, err error) {
	q := `INSERT INTO public.workflow
		(guild_id, starting_discord_events)
	VALUES ($1, $2)
	RETURNING *;`

	err = w.postgresDatabaseService.Get(&workflow, q, guildID, models.DiscordEventTypeArray(startingDiscordEvents))
	return
}

func (w WorkflowRepository) ListByGuildID(guildID string, withStartingNodesIDs bool) (workflows []models.Workflow, err error) {
	q := `SELECT *, snids.starting_nodes_ids FROM public.workflow `

	if withStartingNodesIDs {
		q += `LEFT JOIN LATERAL (
			SELECT ARRAY_AGG(wsn.node_id) AS starting_nodes_ids
			FROM public.workflow_starting_node wsn
			WHERE wsn.workflow_id = id) AS snids ON true `
	}

	q += `WHERE guild_id = $1 AND deleted_at IS NULL`

	err = w.postgresDatabaseService.Db.Select(&workflows, q, guildID)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	// Initialize empty nodes array for each workflow
	for i := range workflows {
		if workflows[i].Nodes == nil {
			workflows[i].Nodes = []models.Node{}
		}
	}

	return
}

// @TODO change this shit
// Use a join to get nodes and use a parameter to decide if we get nodes
// GetByID returns a workflow including its nodes
func (w WorkflowRepository) GetByID(id string) (workflow models.Workflow, err error) {
	q := `SELECT * FROM public.workflow WHERE id = $1 AND deleted_at IS NULL;`
	err = w.postgresDatabaseService.Get(&workflow, q, id)
	if err != nil {
		return
	}

	var nodes []models.Node
	nq := `SELECT * FROM public.node WHERE workflow_id = $1 AND deleted_at IS NULL ORDER BY created_at;`
	err = w.postgresDatabaseService.Db.Select(&nodes, nq, id)
	if err != nil {
		return workflow, oops.Wrap(err)
	}
	workflow.Nodes = nodes
	return
}

func (w WorkflowRepository) GetByStartingDiscordEvents(guildID string, discordEventTypes []models.DiscordEventType) (workflow models.Workflow, err error) {
	q := `SELECT * FROM public.workflow WHERE guild_id = $1 AND starting_discord_events && $2 AND deleted_at IS NULL;`
	err = w.postgresDatabaseService.Get(&workflow, q, guildID, models.DiscordEventTypeArray(discordEventTypes))
	return
}

// Update updates the provided fields of a workflow and returns the updated row
func (w WorkflowRepository) Update(workflowID string, startingDiscordEvents models.DiscordEventTypeArray, startingNodesIDs []string) (updated models.Workflow, err error) {
	q := `UPDATE public.workflow SET starting_discord_events = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *;`
	err = w.postgresDatabaseService.Get(&updated, q, workflowID, startingDiscordEvents)
	if err != nil {
		return models.Workflow{}, err
	}
	// @TODO Change and opitmize this shitty code
	// Update starting nodes IDs
	// First, delete existing starting nodes
	delQ := `DELETE FROM public.workflow_starting_node WHERE workflow_id = $1;`
	err = w.postgresDatabaseService.Exec(delQ, workflowID)
	if err != nil {
		return models.Workflow{}, err
	}
	// Then, insert new starting nodes
	for _, nodeID := range startingNodesIDs {
		insQ := `INSERT INTO public.workflow_starting_node (workflow_id, node_id) VALUES ($1, $2);`
		err = w.postgresDatabaseService.Exec(insQ, workflowID, nodeID)
		if err != nil {
			return models.Workflow{}, err
		}
	}
	return
}

// Delete soft-deletes a workflow
func (w WorkflowRepository) Delete(id string) (err error) {
	q := `UPDATE public.workflow SET deleted_at = NOW() WHERE id = $1;`
	err = w.postgresDatabaseService.Exec(q, id)
	return
}
