package handlers

import (
	slackbot "github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

func Help(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	commands := map[string]string{
		"deploy <VM_NAME> <URI to OVA file> [NETWORK]": "Deploy Virtual Machine from OVA file.",
		"info <VM name>":                               "Information about Virtual Machine.",
		"help":                                         "See the available bot commands.",
	}

	fields := make([]slack.AttachmentField, 0)
	for k, v := range commands {
		fields = append(fields, slack.AttachmentField{
			Title: "@janna " + k,
			Value: v,
		})
	}
	attachment := &slack.Attachment{
		Pretext: "Janna Command List",
		Color:   "#7CD197",
		Fields:  fields,
	}

	// multiple attachments
	attachments := []slack.Attachment{*attachment}

	bot.ReplyWithAttachments(evt, attachments, slackbot.WithTyping)
}
