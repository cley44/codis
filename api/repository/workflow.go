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
func (w WorkflowRepository) Create(startNodeID string, guildID string, startingDiscordEvents []models.DiscordEventType) (workflow models.Workflow, err error) {
	q := `INSERT INTO public.workflow
		(start_node_id, guild_id, starting_discord_events, edges)
	VALUES ($1, $2, $3, $4)
	RETURNING *;`

	// default edges to empty array
	err = w.postgresDatabaseService.Get(&workflow, q, startNodeID, guildID, startingDiscordEvents, []byte("[]"))
	return
}

// List returns workflows optionally filtered by guild id
func (w WorkflowRepository) List(guildID *string) (workflows []models.Workflow, err error) {
	q := `SELECT id, start_node_id, guild_id, starting_discord_events, edges, created_at, updated_at FROM public.workflow WHERE deleted_at IS NULL`
	if guildID != nil {
		err = w.postgresDatabaseService.Db.Select(&workflows, q+" AND guild_id = $1", *guildID)
	} else {
		err = w.postgresDatabaseService.Db.Select(&workflows, q)
	}
	if err != nil {
		return nil, oops.Wrap(err)
	}
	return
}

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

// Update updates the provided fields of a workflow and returns the updated row
func (w WorkflowRepository) Update(workflow models.Workflow) (updated models.Workflow, err error) {
	q := `UPDATE public.workflow SET start_node_id = $1, starting_discord_events = $2, edges = $3, updated_at = NOW() WHERE id = $4 AND deleted_at IS NULL RETURNING *;`
	err = w.postgresDatabaseService.Get(&updated, q, workflow.StartNodeID, workflow.StartingDiscordEvents, workflow.Edges, workflow.ID)
	return
}

// Delete soft-deletes a workflow
func (w WorkflowRepository) Delete(id string) (err error) {
	q := `UPDATE public.workflow SET deleted_at = NOW() WHERE id = $1;`
	err = w.postgresDatabaseService.Exec(q, id)
	return
}
