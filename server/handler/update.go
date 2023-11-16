package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/status"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func MakeUpdate(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	newVersion, err := servertools.RunUpdate(status.ServerVersion, "Festivals-App", "festivals-gateway", "/usr/local/festivals-gateway/update.sh")
	if err != nil {
		log.Error().Err(err).Msg("Failed to update")
		servertools.RespondError(w, http.StatusInternalServerError, "Failed to update")
		return
	}
	servertools.RespondString(w, http.StatusAccepted, newVersion)
}
