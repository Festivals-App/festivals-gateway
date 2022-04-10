package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/status"
	"github.com/Festivals-App/festivals-gateway/server/update"
	"github.com/rs/zerolog/log"
)

func MakeUpdate(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	newVersion, err := update.RunUpdate(status.ServerVersion, "https://github.com/Festivals-App/festivals-gateway/releases/latest", "/usr/local/festivals-gateway/update.sh")
	if err != nil {
		log.Error().Err(err).Msg("Failed to update")
		respondError(w, http.StatusInternalServerError, "Failed to update")
		return
	}

	respondString(w, http.StatusAccepted, newVersion)
}
