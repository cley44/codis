package discord

import (
	"codis/models"
	"codis/utils"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/oauth2"
)

// @TODO Permissions should be passe as parameters
func (svc *DiscordService) GetDiscordInviteLink() string {

	params := oauth2.AuthorizationURLParams{
		RedirectURI: svc.config.Discord.RedirectURI,
		//discord.OAuth2ScopeBot,
		Scopes: []discord.OAuth2Scope{
			discord.OAuth2ScopeIdentify,
			discord.OAuth2ScopeEmail,
			discord.OAuth2ScopeGuilds,
		},
		Permissions: discord.PermissionAdministrator,
	}

	authorizatonURL := oauth2.Client.GenerateAuthorizationURL(svc.oauthClient, params)

	return authorizatonURL
}

func (svc *DiscordService) StartSession(code string, state string) models.DiscordSession {
	oauthSession, _, err := svc.oauthClient.StartSession(code, state)
	if err != nil {
		utils.PrintJSONIndent(err.Error())
		//@TODO should be handled
		panic(err.Error())
	}

	discordSession := models.DiscordSession{
		Session: oauthSession,
	}

	return discordSession
}

func (svc *DiscordService) GetUser(session oauth2.Session) (discordUser *discord.OAuth2User, err error) {
	return svc.oauthClient.GetUser(session)
}
