package discord

import (
	"codis/utils"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/oauth2"
)

// @TODO Permissions should be passe as parameters
func (svc *DiscordService) GetDiscordInviteLink() string {

	params := oauth2.AuthorizationURLParams{
		RedirectURI: svc.config.Discord.RedirectURI,
		Scopes:      []discord.OAuth2Scope{discord.OAuth2ScopeBot, discord.OAuth2ScopeIdentify},
		Permissions: discord.PermissionAdministrator,
	}

	authorizatonURL := oauth2.Client.GenerateAuthorizationURL(svc.oauthClient, params)

	return authorizatonURL
}

func (svc *DiscordService) StartSession(code string, state string) {
	oauthSession, _, err := svc.oauthClient.StartSession(code, state)
	if err != nil {
		//@TODO should be handled
		panic("fuckl")
	}

	oauthUser, err := svc.oauthClient.GetUser(oauthSession)
	if err != nil {
		panic("fuck 2")
	}
	utils.PrintJSONIndent(oauthUser)
}
