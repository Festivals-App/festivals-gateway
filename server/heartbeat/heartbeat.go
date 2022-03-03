package heartbeat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Heartbeat struct {
	Service   string `json:"service"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Available bool   `json:"available"`
}

func SendHeartbeat(url string, serviceKey string, beat *Heartbeat) error {

	heartbeatwave, err := json.Marshal(beat)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(heartbeatwave))
	if err != nil {
		return err
	}

	request.Header.Set("Api-Key", serviceKey)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
