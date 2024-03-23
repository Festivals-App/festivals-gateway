package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/status"
	token "github.com/Festivals-App/festivals-identity-server/jwt"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func GetServices(validator *token.ValidationService, claims *token.UserClaims, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get the service list.")
		servertools.UnauthorizedResponse(w)
		return
	}

	var nodes = []status.MonitorNode{}

	Loadbalancer.ServicesMux.RLock()
	defer Loadbalancer.ServicesMux.RUnlock()

	for _, element := range Loadbalancer.Services {
		for _, service := range element.Services {
			node := status.MonitorNode{Type: service.Name, Location: service.URL.String(), LastSeen: service.SeenAt()}
			nodes = append(nodes, node)
		}
	}
	servertools.RespondJSON(w, http.StatusOK, nodes)
}
