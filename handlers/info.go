package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	slackbot "github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/vterdunov/janna-slack-bot/utils"
	"golang.org/x/net/context"
)

// NetworkInfo provide VM network info
type NetworkInfo struct {
	IP string `json:"guest_ip"`
}

// PowerInfo provide info about VM power state
type PowerInfo struct {
	State string `json:"state"`
}

// InfoVM provide VM info
type InfoVM struct {
	Name         string      `json:"name"`
	UUID         string      `json:"uuid"`
	InstanceUUID string      `json:"instance_uuid"`
	Network      NetworkInfo `json:"network"`
	Power        PowerInfo   `json:"power"`
}

// Info return information about Virtual Machine
func Info(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.WithFields(log.Fields{
		"handler": "Info",
		"message": evt.Msg.Text,
	}).Info("Request for get information about VM")

	jannaAPIAddress := os.Getenv("JANNA_API_ADDRESS")

	var reply string
	errorReply := "Something went wrong."
	vmName := utils.MessageTrim(evt.Msg.Text, "info")[0]
	url := jannaAPIAddress + "/v1/vm?provider_type=vmware&vmname=" + vmName
	resp, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error while request to Janna API")
		bot.Reply(evt, errorReply, false)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"url":           url,
			"responce code": resp.StatusCode,
		}).Error("Responce code not a 200 OK")
		bot.Reply(evt, errorReply, false)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error while reading body")
		bot.Reply(evt, errorReply, false)
		return
	}
	vminfo := InfoVM{}
	err = json.Unmarshal(bodyBytes, &vminfo)
	if err != nil {
		reply = "Error json unmarshal."
		log.WithFields(log.Fields{
			"error": err,
		}).Error(reply)
		bot.Reply(evt, reply, false)
	}

	vmValues := map[string]string{
		"IP address":    vminfo.Network.IP,
		"Power state":   vminfo.Power.State,
		"uuid":          vminfo.UUID,
		"Instance uuid": vminfo.InstanceUUID,
	}

	fields := make([]slack.AttachmentField, 0)
	for k, v := range vmValues {
		fields = append(fields, slack.AttachmentField{
			Title: k,
			Value: v,
		})
	}

	attachment := &slack.Attachment{
		Pretext: vminfo.Name + " Information",
		Color:   "a9a9a9",
		Fields:  fields,
	}

	// multiple attachments
	attachments := []slack.Attachment{*attachment}

	bot.ReplyWithAttachments(evt, attachments, slackbot.WithTyping)
}
