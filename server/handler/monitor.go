package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/status"
)

func GetServices(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	var nodes = []status.MonitorNode{}

	Loadbalancer.ServicesMux.RLock()
	defer Loadbalancer.ServicesMux.RUnlock()

	for _, element := range Loadbalancer.Services {
		for _, service := range element.Services {
			node := status.MonitorNode{Type: service.Name, Location: service.URL.String(), LastSeen: service.SeenAt()}
			nodes = append(nodes, node)
		}
	}
	respondJSON(w, http.StatusOK, nodes)
}
