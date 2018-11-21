// slack package handles incoming and outgoing messages from Slask
package slack

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/internal/bot"
)

var botUserID string
var botMentionTag string

const protocolKey = "slack"

// Run connects to slack API with the provided token
func Run(token string, b *bot.Bot) error {
	api := slack.New(token)
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
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

			botUserID = ev.Info.User.ID
			botMentionTag = fmt.Sprintf("<@%s>", botUserID)
			outgoingChannel := b.RegisterProtocol(protocolKey)
			go sendMessage(rtm, outgoingChannel)

		case *slack.MessageEvent:
			if !isValid(ev) {
				continue
			}

			u, err := api.GetUserInfo(ev.User)
			if err != nil {
				log.Warn().Msg("could not retrieve slack user info")
			}

			msg := bot.MessageData{
				User:     u.Name,
				Message:  ev.Text,
				Protocol: protocolKey,
				Channel:  ev.Channel,
			}

			b.ReceiveMessage(msg)

		case *slack.RTMError:
			log.Error().Str("error", ev.Msg).Int("code", ev.Code).Msg("")

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")

		case *slack.DisconnectedEvent:
			log.Info().Msg("Bot disconnected")

		default:
			// Ignore other events..
		}
	}

	return nil
}

func sendMessage(rtm *slack.RTM, ch chan bot.OutgoingMessage) {
	params := slack.PostMessageParameters{AsUser: true}

	for outMsg := range ch {
		if outMsg.ErrText != "" {
			sendMessageWithError(rtm, outMsg.Channel, outMsg.ErrText)
		}
		fields := make([]slack.AttachmentField, 0)
		fields = append(fields, slack.AttachmentField{
			Value: outMsg.Text,
		})
		attachment := &slack.Attachment{
			Pretext: outMsg.Title,
			Color:   "#7CD197",
			Fields:  fields,
		}
		attachments := []slack.Attachment{*attachment}
		params.Attachments = attachments

		rtm.PostMessage(outMsg.Channel, "", params)
	}
}

// isDirectMessage returns true if this message is in a direct message conversation
func isDirectMessage(ev *slack.MessageEvent) bool {
	return regexp.MustCompile("^D.*").MatchString(ev.Channel)
}

func isValid(ev *slack.MessageEvent) bool {
	if ev.Hidden || botUserID == ev.User {
		return false
	}

	if !strings.HasPrefix(ev.Msg.Text, botMentionTag) && !isDirectMessage(ev) {
		return false
	}

	return true
}

func sendMessageWithError(rtm *slack.RTM, channel, err string) {
	attachment := &slack.Attachment{
		Color: "ff0000",
		Text:  err,
		ID:    1000,
		Title: "Error",
		// Footer: reqID,
	}
	// multiple attachments
	attachments := []slack.Attachment{*attachment}
	params := slack.PostMessageParameters{AsUser: true}
	params.Attachments = attachments

	rtm.PostMessage(channel, "", params)
}
