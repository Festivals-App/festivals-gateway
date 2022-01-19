package handler

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/Festivals-App/festivals-gateway/server/config"
)

func GoToFestivalsAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	url, err := loadbalancedHost("festivals-server")
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	rp.ServeHTTP(w, r)
}

func GoToFestivalsFilesAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	url, err := loadbalancedHost("festivals-fileserver")
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	rp.ServeHTTP(w, r)
}

func loadbalancedHost(serviceIdentifier string) (*url.URL, error) {

	pool, exists := ServicePools[serviceIdentifier]
	if !exists {
		return nil, errors.New("loadbalancer: no available backend server for " + serviceIdentifier)
	}

	peer := pool.GetNextPeer()
	if peer != nil {
		return peer.URL, nil
	}

	return nil, errors.New("loadbalancer: no available backend server for " + serviceIdentifier)
}

func (s *ServicePool) NextIndex() int {
	return int(atomic.AddUint64(&s.Current, uint64(1)) % uint64(len(s.Services)))
}

// GetNextPeer returns next active peer to take a connection
func (s *ServicePool) GetNextPeer() *Service {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	l := len(s.Services) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.Services) // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one
		if s.Services[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.Current, uint64(idx)) // mark the current one
			}
			return s.Services[idx]
		}
	}
	return nil
}
