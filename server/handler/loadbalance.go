package handler

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/loadbalancer"
)

var ServicePools map[string]*loadbalancer.ServicePool = map[string]*loadbalancer.ServicePool{}
var servicePoolMux sync.RWMutex

var Pools sync.Map

func GoToFestivalsWebsite(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	url, err := loadbalancedHost("festivals-website")
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	rp.ServeHTTP(w, r)
}

func GoToFestivalsWebsiteNode(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	url, err := loadbalancedHost("festivals-website-node")
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	rp.ServeHTTP(w, r)
}

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

	servicePoolMux.RLock()
	defer servicePoolMux.RUnlock()

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
