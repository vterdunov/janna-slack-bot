package bot

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/internal/config"
)

// Bot represents a bot
type Bot struct {
	Name            string
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
		JannaAPIAddress: cfg.JannaAPIAddress,
		Logger:          logger,
		Client:          client,
	}
}

// Run bot
func (b *Bot) Run(_ context.Context) error {
	_, err := b.Client.AuthTest()
	if err != nil {
		return errors.Wrap(err, "did not authenticate")
	}

	b.RTM = b.Client.NewRTM()
	go b.RTM.ManageConnection()

	for msg := range b.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.RTMError:
			log.Error().Str("error", ev.Msg).Int("code", ev.Code).Msg("")

		case *slack.ConnectingEvent:
			log.Info().
				Int("connections count", ev.ConnectionCount).
				Int("attempt", ev.Attempt).
				Msg("Connecting bot to Slack")

		case *slack.ConnectedEvent:
			log.Info().
				Str("bot name", ev.Info.User.Name).
				Int("connections count", ev.ConnectionCount).
				Msg("Connected")

			b.UserID = ev.Info.User.ID

		case *slack.MessageEvent:
			b.handleMessageEvent(ev)

		case *slack.DisconnectedEvent:
			log.Info().Msg("Bot disconnected")

		default:
			// Ignore other events..
		}
	}

	return nil
}

func (b *Bot) handleMessageEvent(ev *slack.MessageEvent) {
	// ignore messages from the current user, the bot user
	if b.UserID == ev.User {
		log.Debug().Msg("Ignore messages from the current user, the bot user")
		return
	}

	// check if the message to bot
	botTag := fmt.Sprintf("<@%s>", b.UserID)
	if !strings.HasPrefix(ev.Msg.Text, botTag) && !isDirectMessage(ev) {
		return
	}

	msg := prepareMsg(ev.Text)

	b.routeMessage(msg, ev)
}

// Reply replies to a message event with a simple message.
func (b *Bot) Reply(channel, msg string) {
	b.RTM.SendMessage(b.RTM.NewOutgoingMessage(msg, channel))

}

// ReplyWithAttachments replys to a message event with a Slack Attachments message.
func (b *Bot) ReplyWithAttachments(channel string, attachments []slack.Attachment) {
	params := slack.PostMessageParameters{AsUser: true}
	params.Attachments = attachments

	b.Client.PostMessage(channel, "", params)
}

func (b *Bot) routeMessage(msg string, ev *slack.MessageEvent) {
	//  vm info
	infoRegexp := regexp.MustCompile(`vm\s+info\s+([a-zA-Z0-9-\.]+)`)
	if infoRegexp.MatchString(msg) {
		log.Debug().Msg("calling VM info handler")
		ss := infoRegexp.FindStringSubmatch(msg)
		vmName := ss[1]

		b.vmInfoHandler(ev.Channel, vmName)
		return
	}

	// vm find
	vmFindRegexp := regexp.MustCompile(`vm\s+find\s+([a-zA-Z0-9-\.]+)`)
	if vmFindRegexp.MatchString(msg) {
		log.Debug().Msg("calling VM find handler")
		ss := vmFindRegexp.FindStringSubmatch(msg)
		vmName := ss[1]

		b.vmFindHandler(ev.Channel, vmName)
		return
	}

	// help
	helpRegexp := regexp.MustCompile(`help`)
	if helpRegexp.MatchString(msg) {
		log.Debug().Msg("calling help handler")
		b.helpHandler(ev.Channel)
		return
	}
}

func prepareMsg(text string) string {
	msg := strings.TrimSpace(text)
	return stripDirectMention(msg)
}

// isDirectMessage returns true if this message is in a direct message conversation
func isDirectMessage(ev *slack.MessageEvent) bool {
	return regexp.MustCompile("^D.*").MatchString(ev.Channel)
}

// stripDirectMention removes a leading mention (aka direct mention) from a message string
func stripDirectMention(text string) string {
	r := regexp.MustCompile(`(^<@[a-zA-Z0-9]+>[\:]*[\s]*)?(.*)`)
	return r.FindStringSubmatch(text)[2]
}
