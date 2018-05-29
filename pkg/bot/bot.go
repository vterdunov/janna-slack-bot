package bot

import (
	"context"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/pkg/config"
	"github.com/vterdunov/janna-slack-bot/pkg/helpers"
)

// Bot represents a bot
type Bot struct {
	Name            string
	ChannelID       string
	UserID          string
	JannaAPIAddress string

	Logger *zerolog.Logger
	Client *slack.Client
	RTM    *slack.RTM
}

// New create a new bot instance
func New(cfg *config.Config, client *slack.Client, logger *zerolog.Logger) *Bot {
	return &Bot{
		Name:            cfg.BotName,
		ChannelID:       cfg.ChannelID,
		JannaAPIAddress: cfg.JannaAPIAddress,
		Logger:          logger,
		Client:          client,
	}
}

// Run bot
func (b *Bot) Run(ctx context.Context) error {
	_, err := b.Client.AuthTest()
	if err != nil {
		return errors.Wrap(err, "did not authenticate")
	}

	b.RTM = b.Client.NewRTM()
	go b.RTM.ManageConnection()

	for msg := range b.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.Info().Str("bot_name", ev.Info.User.Name).Msg("Connected")

			b.UserID = ev.Info.User.ID
		case *slack.MessageEvent:
			if err := b.handleMessageEvent(ev); err != nil {
				log.Error().Err(err).Msg("Failed to handle message")
			}
		}
	}

	return nil
}

func (b *Bot) handleMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != b.ChannelID {
		log.Debug().Msg("Only response in specific channel. Ignore.")
		return nil
	}

	// ignore messages from the current user, the bot user
	if b.UserID == ev.User {
		log.Debug().Msg("Ignore messages from the current user, the bot user")
		return nil
	}

	msgs := messageTrim(ev.Msg.Text)

	b.routeMessage(msgs, ev)
	return nil
}

// reply replies to a message event with a simple message.
func reply(rtm *slack.RTM, ev *slack.MessageEvent, msg string) {
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))

}

// replyWithAttachments replys to a message event with a Slack Attachments message.
func replyWithAttachments(c *slack.Client, ev *slack.MessageEvent, attachments []slack.Attachment) {
	params := slack.PostMessageParameters{AsUser: true}
	params.Attachments = attachments

	c.PostMessage(ev.Msg.Channel, "", params)
}

func messageTrim(msg string) []string {
	text := strings.TrimSpace(msg)
	return strings.Fields(text)[1:]
}

func (b *Bot) routeMessage(msgs []string, ev *slack.MessageEvent) {
	switch msgs[0] {
	case "vm":
		switch msgs[1] {
		case "info":
			vmInfoHandler(b.Client, ev, b.JannaAPIAddress, msgs[2], b.RTM)
		}
	default:
		helpHandler(b.Client, ev)
	}
}

func helpHandler(c *slack.Client, ev *slack.MessageEvent) {
	attachments := helpers.HelpAttachments()
	replyWithAttachments(c, ev, attachments)
}

func vmInfoHandler(c *slack.Client, ev *slack.MessageEvent, ja string, vmName string, rtm *slack.RTM) {

	attachments, err := helpers.VMInfo(c, ev, ja, vmName)
	if err != nil {
		log.Error().Err(err).Msg("Could not get VM info")
		reply(rtm, ev, err.Error())
	}
	replyWithAttachments(c, ev, attachments)
}
