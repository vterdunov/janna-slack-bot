package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type list struct {
	List map[string]string `json:"vm_list"`
}

func List(jannaAddr string) ([]string, error) {
	url := fmt.Sprintf("%s/vm", jannaAddr)

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

	vmList := list{}
	err = json.Unmarshal(bodyBytes, &vmList)
	if err != nil {
		return nil, err
	}

	vms := make([]string, len(vmList.List))
	for _, name := range vmList.List {
		vms = append(vms, name)

	}

	return vms, nil
}
