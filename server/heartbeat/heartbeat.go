package heartbeat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Heartbeat struct {
	Service   string `json:"service"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Available bool   `json:"available"`
}

func SendHeartbeat(url string, beat *Heartbeat) {

	heartbeatwave, err := json.Marshal(beat)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(heartbeatwave))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
}
