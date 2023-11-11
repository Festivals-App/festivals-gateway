package handler

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"sync"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/loadbalancer"
)

func GoToFestivalsIdentityAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	goToLoadbalancedHost("festivals-identity-server", conf, w, r)
}

func GoToFestivalsAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	goToLoadbalancedHost("festivals-server", conf, w, r)
}

func GoToFestivalsDatabase(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	goToLoadbalancedHost("festivals-database", conf, w, r)
}

func GoToFestivalsFilesAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	goToLoadbalancedHost("festivals-fileserver", conf, w, r)
}

func GoToFestivalsWebsiteNode(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	goToLoadbalancedHost("festivals-website-node", conf, w, r)
}

func goToLoadbalancedHost(service string, conf *config.Config, w http.ResponseWriter, r *http.Request) {

	host, err := loadbalancedHost(service)
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	rp := httputil.NewSingleHostReverseProxy(host)
	rp.ServeHTTP(w, r)
}

// Allowed load balanced service identifier
var loadBalancedServiceIdentifier = []string{
	"festivals-gateway",
	"festivals-identity-server",
	"festivals-server",
	"festivals-fileserver",
	"festivals-database",
	"festivals-database-node",
	"festivals-website",
	"festivals-website-node",
}

var Loadbalancer = &loadbalancer.LoadBalancer{
	Services:    map[string]*loadbalancer.ServicePool{},
	ServicesMux: sync.RWMutex{},
}

func loadbalancedHost(serviceIdentifier string) (*url.URL, error) {

	if !slices.Contains(loadBalancedServiceIdentifier, serviceIdentifier) {
		return nil, errors.New("loadbalancer: no available backend server for " + serviceIdentifier)
	}

	Loadbalancer.ServicesMux.RLock()
	defer Loadbalancer.ServicesMux.RUnlock()

	services, exists := Loadbalancer.Services[serviceIdentifier]
	if !exists {
		return nil, errors.New("loadbalancer: no available backend server for " + serviceIdentifier)
	}

	service := services.GetNextService()
	if service != nil {
		return service.URL, nil
	}

	return nil, errors.New("loadbalancer: no available backend server for " + serviceIdentifier)
}
