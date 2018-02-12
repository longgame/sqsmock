package worker

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"github.com/greenac/sqsmock/logger"
)

type endpoint string

const (
	workerNewMessageEndpoint endpoint = "/notification-center/new-message"
)

type Interface struct {
	BaseUrl string
}

func (i *Interface) SendNewMessage(pl interface{}) error {
	return i.sendToWorker(pl, workerNewMessageEndpoint)
}

func (i *Interface) sendToWorker(pl interface{}, ep endpoint) error {
	body, err := json.Marshal(pl)
	if err != nil {
		logger.Error("Could not format message to send to worker. Error:", err)
		return err
	}

	// FIXME: Find a way to join these paths
	url := i.BaseUrl + string(ep)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		logger.Error("Failed to create request to:", url, "error:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Could not write message:", pl, "to", url, "error:", err)
		return err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to post message:", pl, "to worker at:", url, "with error:", err)
		return err
	}

	logger.Log("got response sending message to worker url:", url, string(res))
	return nil
}
