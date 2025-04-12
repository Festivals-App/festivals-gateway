package server

import (
	"crypto/tls"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/handler"
	token "github.com/Festivals-App/festivals-identity-server/jwt"
	festivalspki "github.com/Festivals-App/festivals-pki"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	"github.com/rs/zerolog/log"
)

// Server has router and tls configuration
type Server struct {
	Router    *chi.Mux
	Config    *config.Config
	TLSConfig *tls.Config
	Validator *token.ValidationService
}

func NewServer(config *config.Config) *Server {
	server := &Server{}
	server.Initialize(config)
	return server
}

// Initialize the server with predefined configuration
func (s *Server) Initialize(config *config.Config) {

	s.Router = chi.NewRouter()
	s.Config = config

	s.setIdentityService()
	s.setTLSHandling()
	s.setMiddleware()
	s.setRoutes()
}

func (s *Server) setIdentityService() {

	config := s.Config
	val := token.NewValidationService(config.IdentityEndpoint, config.TLSCert, config.TLSKey, config.TLSRootCert, config.ServiceKey, true)
	if val == nil {
		log.Fatal().Msg("failed to create validator")
	}
	s.Validator = val
}

func (s *Server) setTLSHandling() {

	tlsConfig, err := festivalspki.NewServerTLSConfig(s.Config.TLSCert, s.Config.TLSKey, s.Config.TLSRootCert)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set TLS handling")
	}
	s.TLSConfig = tlsConfig
}

func (s *Server) setMiddleware() {
	// tell the ruter which middleware to use
	s.Router.Use(
		// used to log the request
		servertools.Middleware(servertools.TraceLogger(s.Config.TraceLog)),
		// tries to recover after panics
		middleware.Recoverer,
	)
}

// setRouters sets the all required routers
func (s *Server) setRoutes() {

	hr := hostrouter.New()

	base := s.Config.ServiceBindHost + ":" + strconv.Itoa(s.Config.ServicePort)
	if s.Config.ServicePort == 80 || s.Config.ServicePort == 443 {
		base = s.Config.ServiceBindHost
	}

	hr.Map("gateway."+base, GetGatewayRouter(s))
	hr.Map("discovery."+base, GetDiscoveryRouter(s))
	hr.Map("api."+base, GetFestivalsAPIRouter(s))
	hr.Map("database."+base, GetFestivalsDatabaseRouter(s))
	hr.Map("files."+base, GetFestivalsFilesAPIRouter(s))

	// Mount the host router
	s.Router.Mount("/", hr)
}

func (s *Server) Run(conf *config.Config) {

	server := http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              conf.ServiceBindHost + ":" + strconv.Itoa(conf.ServicePort),
		Handler:           s.Router,
		TLSConfig:         s.TLSConfig,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal().Err(err).Str("type", "server").Msg("failed to run server")
	}
}

func GetGatewayRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Get("/health", s.handleRequest(handler.GetHealth))
	r.Get("/version", s.handleRequest(handler.GetVersion))
	r.Get("/info", s.handleRequest(handler.GetInfo))

	r.Get("/log", s.handleRequest(handler.GetLog))
	r.Get("/log/trace", s.handleRequest(handler.GetTraceLog))
	r.Post("/update", s.handleRequest(handler.MakeUpdate))

	return r
}

func GetDiscoveryRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Post("/loversear", s.handleServiceRequest(handler.ReceivedHeartbeat))
	r.Get("/services", s.handleRequest(handler.GetServices))

	return r
}

func GetFestivalsAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.loadbalanceRequest(handler.GoToFestivalsAPI))
	return r
}

func GetFestivalsDatabaseRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.loadbalanceRequest(handler.GoToFestivalsDatabase))
	return r
}

func GetFestivalsFilesAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.loadbalanceRequest(handler.GoToFestivalsFilesAPI))
	return r
}

type RequestHandlerFunction func(config *config.Config, w http.ResponseWriter, r *http.Request)

func (s *Server) loadbalanceRequest(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.Config, w, r)
	})
}

type JWTAuthenticatedHandlerFunction func(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request)

func (s *Server) handleRequest(requestHandler JWTAuthenticatedHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims := token.GetValidClaims(r, s.Validator)
		if claims == nil {
			servertools.UnauthorizedResponse(w)
			return
		}
		requestHandler(s.Validator, claims, w, r)
	})
}

type ServiceKeyAuthenticatedHandlerFunction func(w http.ResponseWriter, r *http.Request)

func (s *Server) handleServiceRequest(requestHandler ServiceKeyAuthenticatedHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		servicekey := token.GetServiceToken(r)
		if servicekey == "" {
			claims := token.GetValidClaims(r, s.Validator)
			if claims != nil && claims.UserRole == token.ADMIN {
				requestHandler(w, r)
				return
			}
			servertools.UnauthorizedResponse(w)
			return
		}
		allServiceKeys := s.Validator.ServiceKeys
		if !slices.Contains(*allServiceKeys, servicekey) {
			servertools.UnauthorizedResponse(w)
			return
		}
		requestHandler(w, r)
	})
}
