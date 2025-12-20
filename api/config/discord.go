package config

type DiscordConfig struct {
	ClientID     string `config:"DISCORD_CLIENT_ID"`
	ClientSecret string `config:"DISCORD_CLIENT_SECRET"`
	RedirectURI  string `config:"DISCORD_REDIRECT_URI"`
	DiscordToken string `config:"DISCORD_TOKEN"`
}
