package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/heartbeat"
)

type Service struct {
	Name         string
	URL          *url.URL
	At           time.Time
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (s *Service) IsAlive() (alive bool) {
	s.mux.RLock()
	alive = s.Alive
	s.mux.RUnlock()
	return
}

type ServicePool struct {
	Services []*Service
	Current  uint64
}

var ServicePools map[string]*ServicePool = map[string]*ServicePool{}

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

	service := &Service{Name: beat.Service, URL: url, At: time.Now(), Alive: beat.Available}
	registerService(service)
	respondCode(w, http.StatusAccepted)
}

func registerService(service *Service) {

	servicePool, exists := ServicePools[service.Name]
	if !exists {
		fmt.Printf("Initially created service pool for: %+v\n", service.Name)
		ServicePools[service.Name] = &ServicePool{Services: []*Service{}, Current: 0}
		registerService(service)
	} else {

		known, index := knownInstance(service, servicePool.Services)
		if known {
			servicePool.Services[index].At = time.Now()
			fmt.Printf("Updated instance: %+v\n", service)

		} else {
			fmt.Printf("Initially added: %+v\n", service)
			servicePool.Services = append(servicePool.Services, service)
		}

		ServicePools[service.Name] = servicePool
	}
}

func knownInstance(service *Service, list []*Service) (bool, int) {

	for i, member := range list {

		a := member.URL
		b := service.URL

		if *a == *b {
			return true, i
		}
	}
	return false, -1
}
