package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/context"

	slackbot "github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
)

var jannaAPIAddress string

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

func main() {
	token, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatal("Provide 'SLACK_TOKEN' environment variable.")
	}

	jannaAPIAddress, ok = os.LookupEnv("JANNA_API_ADDRESS")
	if !ok {
		log.Fatal("Provide 'JANNA_API_ADDRESS' environment variable.")
	}

	bot := slackbot.New(token)

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(^info ).*").MessageHandler(infoHandler)
	toMe.Hear("(^deploy ).*").MessageHandler(deployOVAHandler)
	toMe.Hear("(?i).*").MessageHandler(helpHandler)
	bot.Run()
}

func messageTrim(msg string, botCommand string) []string {
	text := slackbot.StripDirectMention(msg)
	text = strings.TrimPrefix(text, botCommand)
	text = strings.TrimSpace(text)
	return strings.Fields(text)
}

func infoHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.Printf("Slack request, handler: infoHandler, message: %s", evt.Msg.Text)

	var reply string
	vmName := messageTrim(evt.Msg.Text, "info")[0]
	url := jannaAPIAddress + "/v1/vm?provider_type=vmware&vmname=" + vmName
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error while request: %s, err: %s", url, err)
		reply = "Something went wrong."
		bot.Reply(evt, reply, false)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Responce code not a 200 OK, request: %s, responce code: %d", url, resp.StatusCode)
		reply = "Something went wrong."
		bot.Reply(evt, reply, false)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading body, request: %s, err: %s", url, err)
		reply = "Something went wrong."
		bot.Reply(evt, reply, false)
		return
	}
	vminfo := InfoVM{}
	err = json.Unmarshal(bodyBytes, &vminfo)
	if err != nil {
		log.Printf("Error json unmarshal")
	}

	reply = vminfo.Network.IP
	bot.Reply(evt, reply, false)
}

func deployOVAHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.Printf("Slack request, handler: deployHandler, message: %s", evt.Msg.Text)

	var reply string
	msgPayload := messageTrim(evt.Msg.Text, "deploy")
	if len(msgPayload) < 2 {
		reply = "Provide options for deploy command."
		bot.Reply(evt, reply, false)
		return
	}
	vmName := msgPayload[0]
	ovaURL := msgPayload[1]
	ovaURL = strings.TrimPrefix(ovaURL, "<")
	ovaURL = strings.TrimSuffix(ovaURL, ">")
	fmt.Println(ovaURL)
	// if len(msgPayload) > 2 {
	// 	network := msgPayload[2]
	// }
	form := url.Values{
		"provider_type": {"vmware"},
		"message_to":    {"#test-hook"},
		"vmname":        {vmName},
		"ova_url":       {ovaURL},
	}
	resp, err := http.PostForm(jannaAPIAddress+"/v1/vm", form)
	if err != nil {
		log.Printf("Error while request, err: %s", err)
		reply = "Something went wrong."
		bot.Reply(evt, reply, false)
		return
	}
	fmt.Printf("Status code %d", resp.StatusCode)
	reply = resp.Status
	bot.Reply(evt, reply, false)
}

func helpHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
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
