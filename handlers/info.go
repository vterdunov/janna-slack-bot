package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	slackbot "github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
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

func Info(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.Printf("Slack request, handler: Info message: %s", evt.Msg.Text)

	jannaAPIAddress := os.Getenv("JANNA_API_ADDRESS")

	var reply string
	vmName := utils.MessageTrim(evt.Msg.Text, "info")[0]
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
		reply = "Error json unmarshal."
		log.Printf(reply)
		bot.Reply(evt, reply, false)
	}

	reply = vminfo.Network.IP
	bot.Reply(evt, reply, false)
}
