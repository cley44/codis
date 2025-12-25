package handlerAPIDiscord

type DiscordGuildsResponse struct {
	ID            string
	Name          string
	IconURL       *string
	BannerURL     *string
	Owner         bool
	BotInviteLink string
	BotPresent    bool
}
