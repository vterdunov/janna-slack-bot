package bot

import (
	"fmt"

	"github.com/vterdunov/janna-slack-bot/internal/vm"
)

func (b *Bot) vmInfoHandler(msg MessageData) OutgoingMessage {
	om := OutgoingMessage{
		Channel: msg.Channel,
		User:    msg.User,
		Title:   "Virtual Machine Info",
	}

	vmName := msg.Cmd[2]
	vmInfo, err := vm.Info(b.JannaAPIAddress, vmName)
	if err != nil {
		om.ErrText = err.Error()
		return om
	}

	result := fmt.Sprintf(`
IP address: %s
Power state: %s
UUID: %s
`, vmInfo.IP, vmInfo.PowerState, vmInfo.UUID)

	om.Text = result

	return om
}

// func (b *Bot) vmFindHandler(ctx context.Context, channel, pattern string) {
// 	vmList, err := vm.List(b.JannaAPIAddress)
// 	if err != nil {
// 		log.Ctx(ctx).Error().Err(err).Msg("Could not get VM list")
// 		b.ReplyWithError(ctx, channel, "Could not get VM list")
// 		return
// 	}

// 	fields := make([]slack.AttachmentField, 0)

// 	for _, name := range vmList {
// 		if strings.Contains(name, pattern) {
// 			fields = append(fields, slack.AttachmentField{
// 				Value: name,
// 			})
// 		}
// 	}

// 	attachment := &slack.Attachment{
// 		Pretext: "Found Virtual Machines",
// 		Color:   "a9a9a9",
// 		Fields:  fields,
// 	}
// 	// multiple attachments
// 	attachments := []slack.Attachment{*attachment}

// 	b.ReplyWithAttachments(channel, attachments)
// }
