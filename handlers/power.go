package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	slackbot "github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/vterdunov/janna-slack-bot/utils"
	"golang.org/x/net/context"
)

// PowerState provide info about VM power state
type PowerState struct {
	State        string `json:"state"`
	ErrorMessage string `json:"error"`
	Ok           bool   `json:"ok"`
}

// Power return information about Virtual Machine
func Power(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.Printf("Slack request, handler: Power, message: %s", evt.Msg.Text)

	var reply string
	errorReply := "Something went wrong."

	msgPayload := utils.MessageTrim(evt.Msg.Text, "power")
	if len(msgPayload) < 2 {
		bot.Reply(evt, "Provide options for the command.", false)
		return
	}
	vmName := msgPayload[0]
	state := msgPayload[1]

	form := url.Values{
		"provider_type": {"vmware"},
		"message_to":    {"#test-hook"},
		"vmname":        {vmName},
		"state":         {state},
	}

	jannaURL := os.Getenv("JANNA_API_ADDRESS")
	resource := "/v1/vm"
	u, _ := url.ParseRequestURI(jannaURL)
	u.Path = resource
	urlString := u.String()
	payload := bytes.NewBufferString(form.Encode())

	req, _ := http.NewRequest(http.MethodPut, urlString, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error while request, err: %s, request: %s", err, urlString)
		bot.Reply(evt, errorReply, false)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Responce code not a 200 OK, request: %s, responce code: %d", urlString, resp.StatusCode)
		bot.Reply(evt, errorReply, false)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading body, request: %s, err: %s", urlString, err)
		bot.Reply(evt, errorReply, false)
		return
	}
	vmState := PowerState{}
	err = json.Unmarshal(bodyBytes, &vmState)
	if err != nil {
		reply = "Error json unmarshal."
		log.Printf(reply)
		bot.Reply(evt, reply, false)
		return
	}
	if !vmState.Ok {
		bot.Reply(evt, vmState.ErrorMessage, false)
		return
	}
	bot.Reply(evt, "Ok!", false)
}
