package models

type DiscordEventType string

const (
	DiscordEventTypeMessageCreate      DiscordEventType = "discord_event_type_message_create"
	DiscordEventTypeMessageReactionAdd DiscordEventType = "discord_event_type_message_reaction_add"
)
