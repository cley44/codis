package config

type RabbitMQConfig struct {
	Hostname string `config:"RABBITMQ_HOSTNAME"`
	Port     int    `config:"RABBITMQ_PORT"`
	Username string `config:"RABBITMQ_USERNAME"`
	Password string `config:"RABBITMQ_PASSWORD"`
	VHost    string `config:"RABBITMQ_VHOST"`
}

