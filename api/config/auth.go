package config

type AuthConfig struct {
	SessionSecret string `config:"SESSION_SECRET"`
}
