package models

import (
	"time"
)

type Workflow struct {
	ID                    string                `db:"id" json:"id"`
	GuildID               string                `db:"guild_id" json:"guild_id"`
	StartingDiscordEvents DiscordEventTypeArray `db:"starting_discord_events" json:"starting_discord_events"`
	Nodes                 []Node                `db:"nodes" json:"nodes"`

	// Added when we send it into amqp
	// ExecutionID *string `db:"execution_id" json:"execution_id"`

	// List of starting node IDs
	StartingNodesIDs *[]string `db:"starting_nodes_ids" json:"starting_nodes_ids"`

	CreateAt  time.Time  `db:"created_at" json:"created_at"`
	UpdateAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
