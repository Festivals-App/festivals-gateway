package handler

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/config"
)

func GoToFestivalsAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	pool, exists := ServicePools["festivals-server"]
	if !exists {
		respondError(w, http.StatusServiceUnavailable, "No available backend server")
		return
	}

	services := pool.Services

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	service := services[rand.Intn(len(services))]

	url, err := url.Parse(service.URL)
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, "Backend server address not reasonable.")
		return
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	rp.ServeHTTP(w, r)

}

func GoToFestivalsFilesAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	//respondString(w, http.StatusOK, "FESTIVALS_FILES_API")

	u, _ := url.Parse("http://localhost:1910")
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ServeHTTP(w, r)
}
