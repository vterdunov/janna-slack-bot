package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nlopes/slack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/pkg/bot"
	"github.com/vterdunov/janna-slack-bot/pkg/config"
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

	logger.Info().Msg("Starting bot")
	if err := bot.Run(ctx); err != nil {
		fmt.Println("ERROR")
	}
}
