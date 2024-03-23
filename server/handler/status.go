package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/status"
	token "github.com/Festivals-App/festivals-identity-server/jwt"
	servertools "github.com/Festivals-App/festivals-server-tools"
)

func GetVersion(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request) {

	servertools.RespondString(w, http.StatusOK, status.VersionString())
}

func GetInfo(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request) {

	servertools.RespondJSON(w, http.StatusOK, status.InfoString())
}

func GetHealth(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request) {

	servertools.RespondCode(w, status.HealthStatus())
}
