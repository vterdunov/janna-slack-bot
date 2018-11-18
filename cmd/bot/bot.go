package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/internal/bot"
	"github.com/vterdunov/janna-slack-bot/internal/config"
	"github.com/vterdunov/janna-slack-bot/pkg/slack"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to process env var")
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	bot := bot.New(cfg, &logger)

	slack.Run(cfg.BotToken, &bot)
}
