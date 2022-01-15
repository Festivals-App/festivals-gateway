package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/Festivals-App/festivals-gateway/server/config"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func GoToFestivalsAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	//respondString(w, http.StatusOK, "FESTIVALS_API")

	u, _ := url.Parse("http://localhost:10439")
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ServeHTTP(w, r)
}

func GoToFestivalsFilesAPI(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	//respondString(w, http.StatusOK, "FESTIVALS_FILES_API")

	u, _ := url.Parse("http://localhost:1910")
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ServeHTTP(w, r)
}
