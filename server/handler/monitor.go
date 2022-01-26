package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/status"
)

func GetServices(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	var nodes []status.MonitorNode = []status.MonitorNode{}

	servicePoolMux.RLock()

	for _, element := range ServicePools {
		for _, service := range element.Services {
			node := status.MonitorNode{Type: service.Name, Location: service.URL.String(), LastSeen: service.SeenAt()}
			nodes = append(nodes, node)
		}
	}

	servicePoolMux.RUnlock()

	respondJSON(w, http.StatusOK, nodes)
}
