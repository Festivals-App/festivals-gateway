package server

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/handler"
	"github.com/Festivals-App/festivals-gateway/server/logger"
	"github.com/Festivals-App/festivals-identity-server/authentication"
	festivalspki "github.com/Festivals-App/festivals-pki"
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

	s.setTLSHandling()
	s.setMiddleware()
	s.setRoutes()
}

func (s *Server) setTLSHandling() {

	tlsConfig := &tls.Config{}
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	tlsConfig.GetCertificate = festivalspki.LoadServerCertificateHandler(s.Config.TLSCert, s.Config.TLSKey, s.Config.TLSRootCert)
	s.TLSConfig = tlsConfig
}

func (s *Server) setMiddleware() {
	// tell the ruter which middleware to use
	s.Router.Use(
		// used to log the request
		logger.Middleware(logger.TraceLogger("/var/log/festivals-gateway/trace.log")),
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

	hr.Map("website."+base, GetWebsiteNodeRouter(s))
	hr.Map("gateway."+base, GetGatewayRouter(s))
	hr.Map("discovery."+base, GetDiscoveryRouter(s))
	hr.Map("api."+base, GetFestivalsAPIRouter(s))
	hr.Map("database."+base, GetFestivalsDatabaseRouter(s))
	hr.Map("files."+base, GetFestivalsFilesAPIRouter(s))
	hr.Map("identity."+base, GetFestivalsIdentityAPIRouter(s))

	// Mount the host router
	s.Router.Mount("/", hr)
}

func GetWebsiteNodeRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutValidation(handler.GoToFestivalsWebsiteNode))
	return r
}

func GetGatewayRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Get("/health", s.handleRequestWithoutValidation(handler.GetHealth))
	r.Get("/version", s.handleRequestWithoutValidation(handler.GetVersion))
	r.Get("/info", s.handleRequestWithoutValidation(handler.GetInfo))

	r.Get("/log", s.handleAdminRequest(handler.GetLog))
	r.Get("/log/trace", s.handleAdminRequest(handler.GetTraceLog))
	r.Post("/update", s.handleAdminRequest(handler.MakeUpdate))

	return r
}

func GetDiscoveryRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Post("/loversear", s.handleAdminRequest(handler.ReceivedHeartbeat))
	r.Get("/services", s.handleAdminRequest(handler.GetServices))

	return r
}

func GetFestivalsAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutValidation(handler.GoToFestivalsAPI))
	return r
}

func GetFestivalsDatabaseRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutValidation(handler.GoToFestivalsDatabase))
	return r
}

func GetFestivalsFilesAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutValidation(handler.GoToFestivalsFilesAPI))
	return r
}

func GetFestivalsIdentityAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutValidation(handler.GoToFestivalsIdentityAPI))
	return r
}

func (s *Server) Run(conf *config.Config) {

	server := http.Server{
		Addr:      conf.ServiceBindHost + ":" + strconv.Itoa(conf.ServicePort),
		Handler:   s.Router,
		TLSConfig: s.TLSConfig,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		// error handling
		log.Fatal().Err(err).Str("type", "server").Msg("Failed to run server")
	}
}

// function prototype to inject config instance in handleRequest()
type RequestHandlerFunction func(config *config.Config, w http.ResponseWriter, r *http.Request)

func (s *Server) handleAdminRequest(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return authentication.IsEntitled(s.Config.AdminKeys, func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.Config, w, r)
	})
}

func (s *Server) handleRequestWithoutValidation(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.Config, w, r)
	})
}
