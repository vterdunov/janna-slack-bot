package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"
)

type InfoVM struct {
	Summary `json:"summary"`
}

type Summary struct {
	Guest   `json:"Guest"`
	Config  `json:"Config"`
	Runtime `json:"Runtime"`
}

type Guest struct {
	IP string `json:"IpAddress"`
}

type Config struct {
	UUID string `json:"Uuid"`
}

type Runtime struct {
	PowerState string `json:"PowerState"`
}

// VMInfo return information about Virtual Machine
func VMInfo(client *slack.Client, ev *slack.MessageEvent, jannaAPIAddress string, vmName string) ([]slack.Attachment, error) {
	// TODO: get VM UUID
	url := jannaAPIAddress + "/vm/564d2d6d-40fe-7e0a-e871-c4ecb46a19d1"
	// url := jannaAPIAddress + "/vm/" + vmName

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

	vminfo := InfoVM{}
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
