package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ShortInfo collects information about VM
type ShortInfo struct {
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
		return "", errors.New("could not find VM")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	uuid := struct {
		UUID string `json:"uuid"`
	}{}

	if err := json.Unmarshal(body, &uuid); err != nil {
		return "", err
	}

	return uuid.UUID, nil
}

// Info return information about Virtual Machine as Slack attachments
func Info(jannaAddr string, vmName string) (*ShortInfo, error) {
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
		return nil, errors.New("request to Janna API was failed")
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	vminfo := ShortInfo{}
	err = json.Unmarshal(bodyBytes, &vminfo)
	if err != nil {
		return nil, err
	}

	return &vminfo, nil
}
