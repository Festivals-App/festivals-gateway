package loadbalancer

import (
	"net/url"
	"sync"
	"time"
)

type Service struct {
	Name  string
	URL   *url.URL
	At    time.Time
	Alive bool
	mux   sync.RWMutex
}

func (service *Service) IsAlive() (alive bool) {
	service.mux.RLock()
	defer service.mux.RUnlock()
	alive = service.Alive
	return alive
}

func (service *Service) SeenAt() (at time.Time) {
	service.mux.RLock()
	defer service.mux.RUnlock()
	at = service.At
	return at
}

func (service *Service) WasSeen(at time.Time) {
	service.mux.Lock()
	defer service.mux.Unlock()
	service.At = at
}
