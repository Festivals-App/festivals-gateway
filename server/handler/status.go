package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/status"
	servertools "github.com/Festivals-App/festivals-server-tools"
)

func GetVersion(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	servertools.RespondString(w, http.StatusOK, status.VersionString())
}

func GetInfo(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	servertools.RespondJSON(w, http.StatusOK, status.InfoString())
}

func GetHealth(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	servertools.RespondCode(w, status.HealthStatus())
}
