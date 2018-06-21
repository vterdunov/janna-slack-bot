package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"

	"github.com/vterdunov/janna-slack-bot/internal/vm"
)

func (b *Bot) helpHandler(channel string) {
	commands := map[string]string{
		"vm deploy <VM_NAME> <URI to OVA file> [NETWORK]": "Deploy Virtual Machine from OVA file.",
		"vm info <VM name>":                               "Information about Virtual Machine.",
		"vm power <VM name> <on|off|reset|suspend>":       "Change Virtual Machine power state.",
		"vm find <part of full of VMs names>":             "Find VMs by wildcard.",
		"help": "See the available bot commands.",
	}

	fields := make([]slack.AttachmentField, 0)
	for k, v := range commands {
		fields = append(fields, slack.AttachmentField{
			Title: fmt.Sprintf("@%s %s", strings.ToLower(b.Name), k),
			Value: fmt.Sprintf("%s\n", v),
		})
	}

	attachment := &slack.Attachment{
		Pretext: b.Name + " Command List",
		Color:   "#7CD197",
		Fields:  fields,
	}

	// multiple attachments
	attachments := []slack.Attachment{*attachment}

	b.ReplyWithAttachments(channel, attachments)
}

func (b *Bot) vmInfoHandler(ctx context.Context, channel, vmName string) {
	vmInfo, err := vm.Info(b.JannaAPIAddress, vmName)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Could not get VM info")
		b.ReplyWithError(ctx, channel, "Could not get VM info")
		return
	}

	vmValues := map[string]string{
		"IP address":  vmInfo.IP,
		"Power state": vmInfo.PowerState,
		"UUID":        vmInfo.UUID,
	}

	fields := make([]slack.AttachmentField, 0)
	for k, v := range vmValues {
		fields = append(fields, slack.AttachmentField{
			Title: k,
			Value: v,
		})
	}

	attachment := &slack.Attachment{
		Pretext: vmName + " information",
		Color:   "a9a9a9",
		Fields:  fields,
	}

	// multiple attachments
	attachments := []slack.Attachment{*attachment}

	b.ReplyWithAttachments(channel, attachments)
}

func (b *Bot) vmFindHandler(ctx context.Context, channel, pattern string) {
	vmList, err := vm.List(b.JannaAPIAddress)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Could not get VM list")
		b.ReplyWithError(ctx, channel, "Could not get VM list")
		return
	}

	fields := make([]slack.AttachmentField, 0)

	for _, name := range vmList {
		if strings.Contains(name, pattern) {
			fields = append(fields, slack.AttachmentField{
				Value: name,
			})
		}
	}

	attachment := &slack.Attachment{
		Pretext: "Found Virtual Machines",
		Color:   "a9a9a9",
		Fields:  fields,
	}
	// multiple attachments
	attachments := []slack.Attachment{*attachment}

	b.ReplyWithAttachments(channel, attachments)
}
