package utils

import (
	"strings"

	slackbot "github.com/adampointer/go-slackbot"
)

func MessageTrim(msg string, botCommand string) []string {
	text := slackbot.StripDirectMention(msg)
	text = strings.TrimPrefix(text, botCommand)
	text = strings.TrimSpace(text)
	return strings.Fields(text)
}
