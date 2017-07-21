package main

import (
	"os"

	slackbot "github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/vterdunov/janna-slack-bot/handlers"
	"golang.org/x/net/context"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func main() {
	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatal("Provide 'SLACK_TOKEN' environment variable.")
	}

	_, ok = os.LookupEnv("JANNA_API_ADDRESS")
	if !ok {
		log.Fatal("Provide 'JANNA_API_ADDRESS' environment variable.")
	}

	bot := slackbot.New(token)

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(^info ).*").MessageHandler(infoHandler)
	toMe.Hear("(^power ).*").MessageHandler(powerHandler)
	toMe.Hear("(^deploy ).*").MessageHandler(deployOVAHandler)
	toMe.Hear("(?i).*").MessageHandler(helpHandler)
	bot.Run()
}

func infoHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	handlers.Info(ctx, bot, evt)
}

func deployOVAHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	handlers.DeployOVA(ctx, bot, evt)
}

func helpHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	handlers.Help(ctx, bot, evt)
}

func powerHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	handlers.Power(ctx, bot, evt)
}
