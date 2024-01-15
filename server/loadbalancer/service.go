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

var staleServiceTreshhold time.Duration = 13 * time.Second

func (service *Service) IsStale() (alive bool) {
	service.mux.RLock()
	defer service.mux.RUnlock()
	return time.Now().Sub(service.At).Seconds() > staleServiceTreshhold.Seconds()
}

func (service *Service) IsEqualTo(otherService *Service) (alive bool) {
	a := service.URL
	b := otherService.URL
	return *a == *b
}
