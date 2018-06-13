package vm

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// 	"os"

// 	slackbot "github.com/adampointer/go-slackbot"
// 	"github.com/nlopes/slack"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/vterdunov/janna-slack-bot/utils"
// 	"golang.org/x/net/context"
// )

// // PowerState provide info about VM power state
// type PowerState struct {
// 	State        string `json:"state"`
// 	ErrorMessage string `json:"error"`
// 	Ok           bool   `json:"ok"`
// }

// // Power return information about Virtual Machine
// func Power(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
// 	log.WithFields(log.Fields{
// 		"handler": "Power",
// 		"message": evt.Msg.Text,
// 	}).Info("Request for change VM power state")

// 	var reply string
// 	errorReply := "Something went wrong."

// 	msgPayload := utils.MessageTrim(evt.Msg.Text, "power")
// 	if len(msgPayload) < 2 {
// 		bot.Reply(evt, "Provide options for the command.", false)
// 		return
// 	}
// 	vmName := msgPayload[0]
// 	state := msgPayload[1]

// 	form := url.Values{
// 		"provider_type": {"vmware"},
// 		"message_to":    {"#test-hook"},
// 		"vmname":        {vmName},
// 		"state":         {state},
// 	}

// 	jannaURL := os.Getenv("JANNA_API_ADDRESS")
// 	resource := "/v1/vm"
// 	u, _ := url.ParseRequestURI(jannaURL)
// 	u.Path = resource
// 	urlString := u.String()
// 	payload := bytes.NewBufferString(form.Encode())

// 	req, _ := http.NewRequest(http.MethodPut, urlString, payload)
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"error":   err,
// 			"request": urlString,
// 		}).Warn("Error while request")

// 		bot.Reply(evt, errorReply, false)
// 		return
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		log.WithFields(log.Fields{
// 			"response code": resp.StatusCode,
// 		}).Error("Response code not a 200 OK")
// 		bot.Reply(evt, errorReply, false)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"error": err,
// 		}).Error("Error while reading body")
// 		bot.Reply(evt, errorReply, false)
// 		return
// 	}
// 	vmState := PowerState{}
// 	err = json.Unmarshal(bodyBytes, &vmState)
// 	if err != nil {
// 		reply = "Error json unmarshal."
// 		log.WithFields(log.Fields{
// 			"error": err,
// 		}).Error(reply)
// 		bot.Reply(evt, reply, false)
// 		return
// 	}
// 	if !vmState.Ok {
// 		bot.Reply(evt, vmState.ErrorMessage, false)
// 		return
// 	}
// 	bot.Reply(evt, "Ok!", false)
// }
