package discord

import (
	"codis/models"

	"github.com/disgoorg/disgo/discord"
)

func (svc *DiscordService) GetGuildsList(userID string, session models.DiscordSession) (guilds []discord.OAuth2Guild, err error) {
	guilds, err = svc.oauthClient.GetGuilds(session.Session)
	return
}
