package discord

import (
	"codis/models"
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func (svc *DiscordService) GetGuildsList(userID string, session models.DiscordSession) (guilds []discord.OAuth2Guild, err error) {
	guilds, err = svc.oauthClient.GetGuilds(session.Session)
	return
}

func (svc *DiscordService) IsBotAMemberOfguild(guildID snowflake.ID) (guild discord.Guild, exist bool) {
	guild, exist = svc.botClient.Caches().Guild(guildID)
	return
}

func (svc *DiscordService) GetGuildMember(guildID snowflake.ID, userID snowflake.ID) {
	svc.botClient.RequestMembers(context.Background(), guildID, true, "", userID)
}

// func (svc *DiscordService) IsBotAMemberOfguild(guildID snowflake.ID) bool {
// 	_, err := svc.botClient.Rest().GetMember(guildID, snowflake.MustParse("1145992042724466688"))
// 	if err != nil {
// 		return false
// 	}
// 	// utils.PrintJSONIndent(member)
// 	return true
// }
