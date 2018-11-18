package bot

// func (b *Bot) vmInfoHandler(ctx context.Context, channel, vmName string) {
// 	vmInfo, err := vm.Info(b.JannaAPIAddress, vmName)
// 	if err != nil {
// 		log.Ctx(ctx).Error().Err(err).Msg("Could not get VM info")
// 		b.ReplyWithError(ctx, channel, "Could not get VM info")
// 		return
// 	}

// 	vmValues := map[string]string{
// 		"IP address":  vmInfo.IP,
// 		"Power state": vmInfo.PowerState,
// 		"UUID":        vmInfo.UUID,
// 	}

// 	fields := make([]slack.AttachmentField, 0)
// 	for k, v := range vmValues {
// 		fields = append(fields, slack.AttachmentField{
// 			Title: k,
// 			Value: v,
// 		})
// 	}

// 	attachment := &slack.Attachment{
// 		Pretext: vmName + " information",
// 		Color:   "a9a9a9",
// 		Fields:  fields,
// 	}

// 	// multiple attachments
// 	attachments := []slack.Attachment{*attachment}

// 	b.ReplyWithAttachments(channel, attachments)
// }

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
