package heartbeat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

type Heartbeat struct {
	Service   string `json:"service"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Available bool   `json:"available"`
}

func HeartbeatClient(clientCert string, clientKey string) (*http.Client, error) {

	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 20 * time.Second,
	}

	return client, nil
}

func SendHeartbeat(client *http.Client, url string, serviceKey string, beat *Heartbeat) error {

	heartbeatwave, err := json.Marshal(beat)
	if err != nil {
		return err
	}

	/*
		http.Request{
			Method: http.MethodPost,
			URL: url,
			Body: bytes.NewBuffer(heartbeatwave),
		}
	*/
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(heartbeatwave))
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
