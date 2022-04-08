package server

import (
	"net/http"
	"strconv"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/handler"
	"github.com/Festivals-App/festivals-gateway/server/logger"
	"github.com/Festivals-App/festivals-identity-server/authentication"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	"github.com/rs/zerolog/log"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server has router and db instances
type Server struct {
	Router *chi.Mux
	Config *config.Config
}

// Initialize the server with predefined configuration
func (s *Server) Initialize(config *config.Config) {

	s.Router = chi.NewRouter()
	s.Config = config

	s.setMiddleware()
	s.setWalker()
	s.setRoutes()
}

func (s *Server) setMiddleware() {
	// tell the ruter which middleware to use
	s.Router.Use(
		// used to log the request
		logger.Middleware(&log.Logger),
		// tries to recover after panics
		middleware.Recoverer,
	)
}

func (s *Server) setWalker() {

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Info().Msg(method + " " + route + " \n")
		return nil
	}
	if err := chi.Walk(s.Router, walkFunc); err != nil {
		log.Panic().Err(err).Msg("Chi walker walked into error")
	}
}

// setRouters sets the all required routers
func (s *Server) setRoutes() {

	hr := hostrouter.New()

	base := s.Config.ServiceBindHost + ":" + strconv.Itoa(s.Config.ServicePort)

	hr.Map(base, GetGatewayRouter(s))
	hr.Map("discovery."+base, GetDiscoveryRouter(s))
	hr.Map("api."+base, GetFestivalsAPIRouter(s))
	hr.Map("files."+base, GetFestivalsFilesAPIRouter(s))

	// Mount the host router
	s.Router.Mount("/", hr)
}

func GetGatewayRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Get("/health", s.handleRequestWithoutValidation(handler.GetHealth))
	r.Get("/version", s.handleRequestWithoutValidation(handler.GetVersion))
	r.Get("/info", s.handleRequestWithoutValidation(handler.GetInfo))

	r.Get("/log", s.handleAdminRequest(handler.GetLog))
	r.Post("/update", s.handleAdminRequest(handler.MakeUpdate))

	return r
}

func GetDiscoveryRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Post("/loversear", s.handleAdminRequest(handler.ReceivedHeartbeat))
	r.Get("/services", s.handleAdminRequest(handler.GetServices))
	//r.Handle("/monitor", promhttp.Handler())
	return r
}

func GetFestivalsAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutValidation(handler.GoToFestivalsAPI))
	return r
}

func GetFestivalsFilesAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutValidation(handler.GoToFestivalsFilesAPI))
	return r
}

// Run the server on it's router
func (s *Server) Run(host string) {
	//log.Fatal(http.ListenAndServeTLS(host, "/cert", "/keys", s.Router))
	if err := http.ListenAndServe(host, s.Router); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}

// function prototype to inject config instance in handleRequest()
type RequestHandlerFunction func(config *config.Config, w http.ResponseWriter, r *http.Request)

func (s *Server) handleAPIRequest(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return authentication.IsEntitled(s.Config.APIKeys, func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.Config, w, r)
	})
}

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
