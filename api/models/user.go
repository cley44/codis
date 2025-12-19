package models

import (
	_ "codis/config"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/disgoorg/disgo/oauth2"
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Username        string    `db:"username" json:"username"`
	DisplayUsername *string   `db:"display_username" json:"display_username"`
	DiscordID       string    `db:"discord_id"`
	DiscordAvatar   *string   `db:"discord_avatar" json:"discord_avatar"`
	Email           string    `db:"email" json:"email"`

	DiscordSession *DiscordSession `db:"discord_session" json:"-"`

	CreateAt  time.Time  `db:"created_at" json:"created_at"`
	UpdateAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

type DiscordSession struct {
	oauth2.Session
}

// https://golang.org/pkg/database/sql/#Scanner implementation
func (s *DiscordSession) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &s)
	}

	return nil
}

// https://golang.org/pkg/database/sql/driver/#Valuer implementation for postgresql json field
func (s DiscordSession) Value() (driver.Value, error) {
	return json.Marshal(s)
}
