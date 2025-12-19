package handlerAPIDiscord

type DiscordCallbackRequest struct {
	Code    string
	State   string
	GuildID string
}
