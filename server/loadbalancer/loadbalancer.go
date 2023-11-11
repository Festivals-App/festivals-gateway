package loadbalancer

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type LoadBalancer struct {
	Services    map[string]*ServicePool
	ServicesMux sync.RWMutex
}

func (loadbalancer *LoadBalancer) Register(service *Service) {

	loadbalancer.ServicesMux.Lock()
	defer loadbalancer.ServicesMux.Unlock()

	servicePool, exists := loadbalancer.Services[service.Name]
	if !exists {
		// check for service name
		servicePool = &ServicePool{}
		loadbalancer.Services[service.Name] = servicePool
		log.Info().Msg("Discovery service created service pool for: " + service.Name)
	}

	servicePool.UpdateService(service)
}
