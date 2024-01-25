package handler

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/loadbalancer"
	festivalspki "github.com/Festivals-App/festivals-pki"
	servertools "github.com/Festivals-App/festivals-server-tools"
)

//"github.com/rs/zerolog/log"

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
		log.Error().Err(err).Msg("Faile to get loadbalances host.")
		servertools.RespondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	tls, err := serverTLSConfig(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Faile to load server TLS configuration.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(host)
	reverseProxy.Transport = &http.Transport{
		TLSClientConfig: tls,
	}
	reverseProxy.ServeHTTP(w, r)
}

func serverTLSConfig(conf *config.Config) (*tls.Config, error) {

	rootCertPool, err := festivalspki.LoadCertificatePool(conf.TLSRootCert)
	if err != nil {
		log.Fatal().Err(err).Msg("Faile to create certificate pool with root CA.")
		return nil, err
	}

	certs, err := tls.LoadX509KeyPair(conf.TLSCert, conf.TLSKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Faile to load server certificate and key.")
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: []tls.Certificate{certs},
	}

	return tlsConfig, nil
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

	Loadbalancer.ServicesMux.Lock()
	defer Loadbalancer.ServicesMux.Unlock()

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
