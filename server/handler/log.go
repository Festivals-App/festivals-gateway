package handler

import (
	"errors"
	"net/http"
	"os"

	token "github.com/Festivals-App/festivals-identity-server/jwt"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func GetLog(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get server log.")
		servertools.UnauthorizedResponse(w)
		return
	}
	l, err := Log("/var/log/festivals-gateway/info.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get info log.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondString(w, http.StatusOK, l)
}

func GetTraceLog(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get server trace log.")
		servertools.UnauthorizedResponse(w)
		return
	}
	l, err := Log("/var/log/festivals-gateway/trace.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get trace log.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondString(w, http.StatusOK, l)
}

func Log(location string) (string, error) {

	l, err := os.ReadFile(location)
	if err != nil {
		return "", errors.New("Failed to read log file at: '" + location + "' with error: " + err.Error())
	}
	return string(l), nil
}
