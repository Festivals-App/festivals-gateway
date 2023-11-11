package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/heartbeat"
	"github.com/Festivals-App/festivals-gateway/server/loadbalancer"

	"github.com/rs/zerolog/log"
)

func ReceivedHeartbeat(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	var beat heartbeat.Heartbeat
	err := json.NewDecoder(r.Body).Decode(&beat)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode send heartbeat payload")
		respondError(w, http.StatusBadRequest, "i don't have any feelings for you")
		return
	}

	url, err := url.Parse(beat.Host + ":" + strconv.Itoa(beat.Port))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse heartbeat sender url")
		respondError(w, http.StatusBadRequest, "i don't have any feelings for you")
		return
	}

	service := &loadbalancer.Service{Name: beat.Service, URL: url, At: time.Now(), Alive: beat.Available}
	go registerService(service)
	respondCode(w, http.StatusAccepted)
}

func registerService(service *loadbalancer.Service) {

	Loadbalancer.ServicesMux.Lock()
	defer Loadbalancer.ServicesMux.Unlock()

	servicePool, exists := Loadbalancer.Services[service.Name]
	if !exists {
		// check for service name
		servicePool = &loadbalancer.ServicePool{}
		Loadbalancer.Services[service.Name] = servicePool
		log.Info().Msg("Discovery service created service pool for: " + service.Name)
	}

	servicePool.UpdateService(service)
}
