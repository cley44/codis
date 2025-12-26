package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/lib/pq"
)

type DiscordEventType string

const (
	DiscordEventTypeMessageCreate      DiscordEventType = "discord_event_type_message_create"
	DiscordEventTypeMessageReactionAdd DiscordEventType = "discord_event_type_message_reaction_add"
)

// DiscordEventTypeArray is a custom type for PostgreSQL array of DiscordEventType
type DiscordEventTypeArray []DiscordEventType

// Scan implements sql.Scanner for DiscordEventTypeArray
func (a *DiscordEventTypeArray) Scan(src interface{}) error {
	var strArray pq.StringArray
	if err := strArray.Scan(src); err != nil {
		return err
	}

	*a = make(DiscordEventTypeArray, len(strArray))
	for i, s := range strArray {
		(*a)[i] = DiscordEventType(s)
	}
	return nil
}

// Value implements driver.Valuer for DiscordEventTypeArray
func (a DiscordEventTypeArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	strArray := make(pq.StringArray, len(a))
	for i, v := range a {
		strArray[i] = string(v)
	}
	return strArray.Value()
}

// MarshalJSON implements json.Marshaler
func (a DiscordEventTypeArray) MarshalJSON() ([]byte, error) {
	if a == nil {
		return json.Marshal([]DiscordEventType{})
	}
	return json.Marshal([]DiscordEventType(a))
}

// UnmarshalJSON implements json.Unmarshaler
func (a *DiscordEventTypeArray) UnmarshalJSON(data []byte) error {
	var arr []DiscordEventType
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*a = DiscordEventTypeArray(arr)
	return nil
}
