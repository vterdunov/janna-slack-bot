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

	// Message is a raw message itself
	Message string

	// Protocol show which service send the message
	Protocol string

	Channel string

	// Cmd parsed message to separated words
	Cmd []string
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
		msg.Cmd = prepareMsg(msg.Message)
		b.routeMessage(msg)
	}
}

func (b *Bot) routeMessage(msg MessageData) error {
	if msg.Cmd[0] == "help" {
		om := helpHandler(msg)
		b.OutgoingMessages[msg.Protocol] <- om
	}

	if len(msg.Cmd) < 2 {
		om := helpHandler(msg)
		om.Title = "Unknown command"
		b.OutgoingMessages[msg.Protocol] <- om
	}

	switch msg.Cmd[1] {
	case "vm":
		switch msg.Cmd[0] {
		case "get":
			om := b.vmInfoHandler(msg)
			b.OutgoingMessages[msg.Protocol] <- om
		case "create":
			fmt.Println("create VM")
		case "delete":
			fmt.Println("delete VM")
		}

	case "snapshot", "snapshots":
		switch msg.Cmd[0] {
		case "get":
			fmt.Println("get VM shapshot")
		case "create":
			fmt.Println("create VM shapshot")
		case "delete":
			fmt.Println("delete VM shapshot")
		}

	case "screenshot":
		switch msg.Cmd[0] {
		case "get":
			fmt.Println("get VM screenshot")
		}

	case "power":
		switch msg.Cmd[0] {
		case "get":
			fmt.Println("get VM power")
		case "create":
			fmt.Println("create VM power")
		default:
			fmt.Println("Unknown power command")
		}

	default:
		fmt.Println("unknown sub-command")
	}

	return nil
}

func prepareMsg(text string) []string {
	msg := strings.TrimSpace(text)
	msg = stripDirectMention(msg)
	cmd := strings.Split(msg, " ")
	return cmd
}

func helpHandler(msg MessageData) OutgoingMessage {
	commands := `
*get vm <name>*
Find and return short information about the Virtual Machine.
*delete vm <name>*
Delete the Virtual Machine

*get snapshot[s] <name>*
List the Virtual Machine snapshots
*create vm snapshot <name>*
*get snapshot[s] <name> <snapshot name>*
Create snapshot for the Virtual Machine
*delete snapshot <name> <snapshot name>*
Delete the snapshot
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
