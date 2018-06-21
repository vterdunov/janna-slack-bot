package main

import (
	"context"
	"os"

	"github.com/nlopes/slack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/internal/bot"
	"github.com/vterdunov/janna-slack-bot/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to process env var")
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	client := slack.New(cfg.BotToken)

	ctx := context.Background()

	bot := bot.New(cfg, client, &logger)

	logger.Info().Msg("Run bot")
	if err := bot.Run(ctx); err != nil {
		logger.Error().Err(err).Msg("error while bot is rinning")
	}
}
