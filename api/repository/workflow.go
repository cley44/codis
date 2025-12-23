package repository

import (
	"codis/models"

	"github.com/samber/do/v2"
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
		(start_node_id, guild_id, starting_discord_events)
	VALUES ($1, $2, $3)
	RETURNING *;`

	err = w.postgresDatabaseService.Get(&workflow, q, startNodeID, guildID, startingDiscordEvents)
	return
}
