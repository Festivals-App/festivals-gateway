package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/heartbeat"
	"github.com/Festivals-App/festivals-gateway/server/loadbalancer"
)

func ReceivedHeartbeat(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	var beat heartbeat.Heartbeat
	err := json.NewDecoder(r.Body).Decode(&beat)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read heartbeat!")
		return
	}

	url, err := url.Parse("http://" + beat.Host + ":" + strconv.Itoa(beat.Port))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read heartbeat URL!")
		return
	}

	service := &loadbalancer.Service{Name: beat.Service, URL: url, At: time.Now(), Alive: beat.Available}
	go registerService(service)
	respondCode(w, http.StatusAccepted)
}

func registerService(service *loadbalancer.Service) {

	servicePoolMux.Lock()

	servicePool, exists := ServicePools[service.Name]

	if !exists {

		fmt.Printf("Initially created service pool for: %+v\n", service.Name)
		ServicePools[service.Name] = &loadbalancer.ServicePool{Services: []*loadbalancer.Service{}, Current: 0}

		servicePoolMux.Unlock()

		registerService(service)
	} else {

		servicePoolMux.Unlock()

		servicePool.UpdateService(service)
	}
}
