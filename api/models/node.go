package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Node struct {
	ID         string          `db:"id" json:"id"`
	WorkflowID string          `db:"workflow_id" json:"workflow_id"`
	Type       DiscordNodeType `db:"type" json:"type"`
	NextNodeID *string         `db:"next_node_id" json:"next_node_id"`

	Data NodeData `db:"data" json:"data"`

	// Added when we send it into amqp
	ExecutionID *string `db:"execution_id" json:"execution_id"`

	CreateAt  time.Time  `db:"created_at" json:"created_at"`
	UpdateAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

type NodeData struct {
	RoleID         *string `json:"role_id,omitempty"`
	ChannelID      *string `json:"channel_id,omitempty"`
	MessageContent *string `json:"message_content,omitempty"`
}

func (s *NodeData) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &s)
	}

	return nil
}

// https://golang.org/pkg/database/sql/driver/#Valuer implementation for postgresql json field
func (s NodeData) Value() (driver.Value, error) {
	return json.Marshal(s)
}
