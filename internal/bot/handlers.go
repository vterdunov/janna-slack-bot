package bot

import (
	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/internal/vm"
)

func (b *Bot) helpHandler(ev *slack.MessageEvent) {
	commands := map[string]string{
		"vm deploy <VM_NAME> <URI to OVA file> [NETWORK]": "Deploy Virtual Machine from OVA file.",
		"vm info <VM name>":                               "Information about Virtual Machine.",
		"vm power <VM name> <on|off|reset|suspend>":       "Change Virtual Machine power state.",
		"help": "See the available bot commands.",
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

	b.ReplyWithAttachments(ev.Channel, attachments)
}

func (b *Bot) vmInfoHandler(ev *slack.MessageEvent, jannaAddr string, vmName string) {
	attachments, err := vm.Info(jannaAddr, vmName)
	if err != nil {
		log.Error().Err(err).Msg("Could not get VM info")
		b.Reply(ev.Channel, err.Error())
	}

	b.ReplyWithAttachments(ev.Channel, attachments)
}
