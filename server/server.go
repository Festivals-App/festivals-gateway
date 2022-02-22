package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/handler"
	"github.com/Festivals-App/festivals-identity-server/authentication"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server has router and db instances
type Server struct {
	Router *chi.Mux
	Config *config.Config
}

// Initialize the server with predefined configuration
func (s *Server) Initialize(config *config.Config) {

	// create router
	s.Router = chi.NewRouter()

	// set config
	s.Config = config

	// prepare server
	s.setMiddleware()
	s.setWalker()
	s.setRoutes()
}

func (s *Server) setMiddleware() {
	// tell the ruter which middleware to use
	s.Router.Use(
		// used to log the request to the console | development
		//middleware.Logger,
		// tries to recover after panics
		middleware.Recoverer,
	)
}

func (s *Server) setWalker() {

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s \n", method, route)
		return nil
	}
	if err := chi.Walk(s.Router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
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
	r.Get("/health", s.handleRequestWithoutAuthentication(handler.GetHealth))
	r.Get("/version", s.handleRequestWithoutAuthentication(handler.GetVersion))
	r.Get("/info", s.handleRequestWithoutAuthentication(handler.GetInfo))
	return r
}

func GetDiscoveryRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Post("/loversear", s.handleRequestWithoutAuthentication(handler.ReceivedHeartbeat))
	r.Get("/services", s.handleRequestWithoutAuthentication(handler.GetServices))
	//r.Handle("/monitor", promhttp.Handler())
	return r
}

func GetFestivalsAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutAuthentication(handler.GoToFestivalsAPI))
	return r
}

func GetFestivalsFilesAPIRouter(s *Server) chi.Router {

	r := chi.NewRouter()
	r.Handle("/*", s.handleRequestWithoutAuthentication(handler.GoToFestivalsFilesAPI))
	return r
}

// Run the server on it's router
func (s *Server) Run(host string) {
	//log.Fatal(http.ListenAndServeTLS(host, "/cert", "/keys", s.Router))
	log.Fatal(http.ListenAndServe(host, s.Router))
}

// function prototype to inject config instance in handleRequest()
type RequestHandlerFunction func(config *config.Config, w http.ResponseWriter, r *http.Request)

// inject Config in handler functions
func (s *Server) handleRequest(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return authentication.IsAuthenticated(s.Config.APIKeys, func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.Config, w, r)
	})
}

// inject Config in handler functions
func (s *Server) handleRequestWithoutAuthentication(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.Config, w, r)
	})
}
