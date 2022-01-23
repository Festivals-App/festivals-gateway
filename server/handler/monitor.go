package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/status"
)

func GetServices(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	var nodes []status.MonitorNode = []status.MonitorNode{}

	for _, element := range ServicePools {
		for _, service := range element.Services {
			nodes = append(nodes, service.monitorNode())
		}
	}

	respondJSON(w, http.StatusOK, nodes)
}

func (s *Service) monitorNode() status.MonitorNode {
	return status.MonitorNode{Type: s.Name, Location: s.URL.String(), LastSeen: s.At}
}
