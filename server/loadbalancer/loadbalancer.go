package loadbalancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type ServicePool struct {
	Services []*Service
	Current  uint64
	mux      sync.RWMutex
}

func (sp *ServicePool) GetNextPeer() *Service {

	sp.mux.Lock()
	defer sp.mux.Unlock()

	// loop entire backends to find out an Alive backend
	next := sp.nextIndex()

	numberOfServices := len(sp.Services)
	l := numberOfServices + next // start from next and move a full cycle

	for i := next; i < l; i++ {
		idx := i % numberOfServices // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one
		if sp.Services[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&sp.Current, uint64(idx)) // mark the current one
			}
			return sp.Services[idx]
		}
	}
	return nil
}

func (sp *ServicePool) nextIndex() int {
	return int(atomic.AddUint64(&sp.Current, uint64(1)) % uint64(len(sp.Services)))
}

func (sp *ServicePool) UpdateService(service *Service) {

	sp.mux.Lock()
	defer sp.mux.Unlock()

	known, index := sp.knownInstance(service)

	if known {
		sp.Services[index].WasSeen(time.Now())
	} else {
		sp.Services = append(sp.Services, service)
	}
}

func (sp *ServicePool) knownInstance(service *Service) (bool, int) {

	for i, member := range sp.Services {

		a := member.URL
		b := service.URL

		if *a == *b {
			return true, i
		}
	}
	return false, -1
}

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

func (s *Service) SeenAt() (at time.Time) {
	s.mux.RLock()
	at = s.At
	s.mux.RUnlock()
	return
}

func (s *Service) WasSeen(at time.Time) {
	s.mux.Lock()
	s.At = at
	s.mux.Unlock()
}
