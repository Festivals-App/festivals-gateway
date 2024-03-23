package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/status"
	token "github.com/Festivals-App/festivals-identity-server/jwt"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func MakeUpdate(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to update server.")
		servertools.UnauthorizedResponse(w)
		return
	}
	newVersion, err := servertools.RunUpdate(status.ServerVersion, "Festivals-App", "festivals-gateway", "/usr/local/festivals-gateway/update.sh")
	if err != nil {
		log.Error().Err(err).Msg("Failed to update")
		servertools.RespondError(w, http.StatusInternalServerError, "Failed to update")
		return
	}
	servertools.RespondString(w, http.StatusAccepted, newVersion)
}
