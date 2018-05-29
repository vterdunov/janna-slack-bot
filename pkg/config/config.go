package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config is the bot config
type Config struct {
	BotToken        string `envconfig:"BOT_TOKEN" required:"true"`
	JannaAPIAddress string `envconfig:"JANNA_API_ADDRESS" required:"true"`
	BotName         string `default:"Janna"`
	ChannelID       string `envconfig:"CHANNEL_ID" required:"true"`
}

// Load config from environment
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
