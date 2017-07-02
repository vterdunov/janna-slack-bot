package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	slackbot "github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/vterdunov/janna-slack-bot/utils"
	"golang.org/x/net/context"
)

// DeployOVA deploy Virtual Machine from OVA file
func DeployOVA(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.Printf("Slack request, handler: deployHandler, message: %s", evt.Msg.Text)

	var reply string
	msgPayload := utils.MessageTrim(evt.Msg.Text, "deploy")
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
	jannaAPIAddress := os.Getenv("JANNA_API_ADDRESS")
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
