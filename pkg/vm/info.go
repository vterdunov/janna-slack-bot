package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"
)

type info struct {
	summary `json:"summary"`
}

type summary struct {
	guest   `json:"Guest"`
	config  `json:"Config"`
	runtime `json:"Runtime"`
}

type guest struct {
	IP string `json:"IpAddress"`
}

type config struct {
	UUID string `json:"Uuid"`
}

type runtime struct {
	PowerState string `json:"PowerState"`
}

func uuidByName(apiAddr string, vmName string) (string, error) {
	url := fmt.Sprintf("%s/find/vm?path=%s", apiAddr, vmName)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Could not get Virtual Machine UUID")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(body)

	// FIX: change the stub
	return "564d2d6d-40fe-7e0a-e871-c4ecb46a19d1", nil
}

// Info return information about Virtual Machine as Slack attachments
func Info(jannaAddr string, vmName string) ([]slack.Attachment, error) {
	uuid, err := uuidByName(jannaAddr, vmName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/vm/%s", jannaAddr, uuid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Request to Janna API was failed. Response code is not 200 OK")
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	vminfo := info{}
	err = json.Unmarshal(bodyBytes, &vminfo)
	if err != nil {
		return nil, err
	}

	vmValues := map[string]string{
		"IP address":  vminfo.IP,
		"Power state": vminfo.PowerState,
		"UUID":        vminfo.UUID,
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

	return attachments, nil
}
