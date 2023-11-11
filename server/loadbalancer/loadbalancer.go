package loadbalancer

import "sync"

type LoadBalancer struct {
	Services    map[string]*ServicePool
	ServicesMux sync.RWMutex
}
