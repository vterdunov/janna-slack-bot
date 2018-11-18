package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog"

	"github.com/vterdunov/janna-slack-bot/internal/config"
)

// Bot represents a bot
type Bot struct {
	Name            string
	JannaAPIAddress string
	Logger          *zerolog.Logger

	Protocols        []string
	Messages         chan MessageData
	OutgoingMessages map[string]chan MessageData
}

type MessageData struct {
	// User who sent the message
	User string

	// Message itself
	Message string

	// Protocol show which service send the message
	Protocol string

	Channel string
}

type contextKey int

func (c contextKey) String() string {
	return string(c)
}

const requestID contextKey = iota

// New create a new bot instance
func New(cfg *config.Config, logger *zerolog.Logger) Bot {
	msgs := make(chan MessageData, 500)
	b := Bot{
		Name:            cfg.BotName,
		JannaAPIAddress: cfg.JannaAPIAddress,
		Logger:          logger,
		Messages:        msgs,
	}

	go b.handleMessages()

	return b
}

// RegisterProtocol register a protocol for send an answer
func (b *Bot) RegisterProtocol(protocol string) chan MessageData {
	b.Protocols = append(b.Protocols, protocol)

}

// ReceiveMessage must be called by a protocol upon receiving a message
func (b *Bot) ReceiveMessage(msg MessageData) {
	b.Messages <- msg
}

func (b *Bot) handleMessages() {
	for msg := range b.Messages {
		prepareMsg(msg.Message)
		b.routeMessage(msg)
	}
}

// // Reply replies to a message event with a simple message.
// func (b *Bot) Reply(channel, msg string) {
// 	b.RTM.SendMessage(b.RTM.NewOutgoingMessage(msg, channel))

// }

// // ReplyWithAttachments replys to a message event with a Slack Attachments message.
// func (b *Bot) ReplyWithAttachments(channel string, attachments []slack.Attachment) {
// 	params := slack.PostMessageParameters{AsUser: true}
// 	params.Attachments = attachments

// 	b.Client.PostMessage(channel, "", params)
// }

// // ReplyWithError replys to a message event with an error message.
// func (b *Bot) ReplyWithError(ctx context.Context, channel, err string) {
// 	reqID, ok := ctx.Value(requestID).(string)
// 	if !ok {
// 		log.Ctx(ctx).Error().Msg("Could not get request ID")
// 	}

// 	attachment := &slack.Attachment{
// 		Color:  "ff0000",
// 		Text:   err,
// 		ID:     1000,
// 		Title:  "Error",
// 		Footer: reqID,
// 	}
// 	// multiple attachments
// 	attachments := []slack.Attachment{*attachment}
// 	params := slack.PostMessageParameters{AsUser: true}
// 	params.Attachments = attachments

// 	b.Client.PostMessage(channel, "", params)
// }

func (b *Bot) routeMessage(msg MessageData) error {
	switch msg.Message {
	case "help":
		fmt.Println("This routed to help handler")
	default:
		fmt.Println("Unknow command")
	}
	// //  vm info
	// infoRegexp := regexp.MustCompile(`vm\s+info\s+([a-zA-Z0-9-\.]+)`)
	// if infoRegexp.MatchString(msg) {
	// 	log.Ctx(ctx).Debug().Msg("calling VM info handler")
	// 	ss := infoRegexp.FindStringSubmatch(msg)
	// 	vmName := ss[1]

	// 	b.vmInfoHandler(ctx, ev.Channel, vmName)
	// 	return
	// }

	// // vm find
	// vmFindRegexp := regexp.MustCompile(`vm\s+find\s+([a-zA-Z0-9-\.]+)`)
	// if vmFindRegexp.MatchString(msg) {
	// 	log.Ctx(ctx).Debug().Msg("calling VM find handler")
	// 	ss := vmFindRegexp.FindStringSubmatch(msg)
	// 	vmName := ss[1]

	// 	b.vmFindHandler(ctx, ev.Channel, vmName)
	// 	return
	// }

	// // help
	// helpRegexp := regexp.MustCompile(`help`)
	// if helpRegexp.MatchString(msg) {
	// 	log.Ctx(ctx).Debug().Msg("calling help handler")
	// 	b.helpHandler(ev.Channel)
	// 	return
	// }
	return nil
}

func prepareMsg(text string) string {
	msg := strings.TrimSpace(text)
	return stripDirectMention(msg)
}

// // isDirectMessage returns true if this message is in a direct message conversation
// func isDirectMessage(ev *slack.MessageEvent) bool {
// 	return regexp.MustCompile("^D.*").MatchString(ev.Channel)
// }

// stripDirectMention removes a leading mention (aka direct mention) from a message string
func stripDirectMention(text string) string {
	r := regexp.MustCompile(`(^<@[a-zA-Z0-9]+>[\:]*[\s]*)?(.*)`)
	return r.FindStringSubmatch(text)[2]
}
