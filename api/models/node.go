package models

import "time"

type Node struct {
	ID         string          `db:"id" json:"id"`
	WorkflowID string          `db:"workflow_id" json:"workflow_id"`
	Type       DiscordNodeType `db:"type" json:"type"`
	NextNodeID *string         `db:"next_node_id" json:"next_node_id"`

	// Added when we send it into amqp
	ExecutionID *string `db:"execution_id" json:"execution_id"`

	CreateAt  time.Time  `db:"created_at" json:"created_at"`
	UpdateAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
