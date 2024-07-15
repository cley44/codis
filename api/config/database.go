package config

type PostgresConfig struct {
	URI      string
	Hostname string `config:"POSTGRES_HOSTNAME"`
	Port     int    `config:"POSTGRES_PORT"`
	Username string `config:"POSTGRES_USERNAME"`
	Password string `config:"POSTGRES_PASSWORD"`
	Database string `config:"POSTGRES_DATABASE"`
}
