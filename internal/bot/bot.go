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
	OutgoingMessages map[string]chan OutgoingMessage
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

type OutgoingMessage struct {
	User    string
	Title   string
	Text    string
	ErrText string
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
		Name:             cfg.BotName,
		JannaAPIAddress:  cfg.JannaAPIAddress,
		Logger:           logger,
		Messages:         msgs,
		OutgoingMessages: map[string]chan OutgoingMessage{},
	}

	go b.handleMessages()

	return b
}

// RegisterProtocol register a protocol for send an answer
func (b *Bot) RegisterProtocol(protocol string) chan OutgoingMessage {
	b.Protocols = append(b.Protocols, protocol)
	omChan := make(chan OutgoingMessage, 100)
	b.OutgoingMessages[protocol] = omChan

	return omChan
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

func (b *Bot) routeMessage(msg MessageData) error {
	switch msg.Message {
	case "help":
		om := helpHandler(msg)
		b.OutgoingMessages[msg.Protocol] <- om
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

func helpHandler(msg MessageData) OutgoingMessage {
	commands := `
*vm deploy <VM_NAME> <URI to OVA file> [NETWORK]*
Deploy Virtual Machine from OVA file

*vm info <VM name>*
Information about Virtual Machine

*vm power <VM name> <on|off|reset|suspend>*
Change Virtual Machine power state

*vm find <part of full of VMs names>*
Find VMs by wildcard
`

	return OutgoingMessage{
		Channel: msg.Channel,
		User:    msg.User,
		Title:   "Available bot commands",
		Text:    commands,
	}
}

// stripDirectMention removes a leading mention (aka direct mention) from a message string
func stripDirectMention(text string) string {
	r := regexp.MustCompile(`(^<@[a-zA-Z0-9]+>[\:]*[\s]*)?(.*)`)
	return r.FindStringSubmatch(text)[2]
}
