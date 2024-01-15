package loadbalancer

import (
	"sync"
	"sync/atomic"
	"time"
)

type ServicePool struct {
	Services []*Service
	Current  uint64
	mux      sync.RWMutex
}

// Loop over entire backends to find out an Alive and active backend service
func (servicePool *ServicePool) GetNextService() *Service {

	servicePool.mux.Lock()
	defer servicePool.mux.Unlock()

	numberOfServices := uint64(len(servicePool.Services))
	if numberOfServices == 0 {
		return nil
	}

	next := servicePool.Current + 1
	if next >= numberOfServices { // out of bound index
		next = 0
	}

	hasStaleService := false
	didRoundtrip := false

	for idx := next; idx < numberOfServices; idx++ {

		// if we have an alive backend, use it and store if its not the original one
		service := servicePool.Services[idx]
		isAlive := service.IsAlive()
		isStale := service.IsStale()

		if isAlive {
			if !isStale {
				servicePool.Current = next
				return service
			} else {
				hasStaleService = true
			}
		}

		// we already iterated from next to numberOfServices
		// so we break the loop if we went from 0 to next
		if didRoundtrip {
			if idx == next {
				break
			}
		}

		// we iterated from next to numberOfServices and reached the last index,
		// now we want to iterate from 0 to next
		if idx+1 == numberOfServices {
			if !didRoundtrip {
				idx = 0
				didRoundtrip = true
			}
		}
	}

	if hasStaleService {

		updatedServices := []*Service{}
		for i := range servicePool.Services {
			service := servicePool.Services[i]
			isStale := service.IsStale()
			if !isStale {
				updatedServices = append(updatedServices, service)
			}
		}
		servicePool.Services = updatedServices
		servicePool.Current = 0
	}

	return nil
}

func (servicePool *ServicePool) NextIndex() int {

	atomic.AddUint64(&servicePool.Current, uint64(1))
	numberOfServices := uint64(len(servicePool.Services))
	return int(servicePool.Current % numberOfServices)
}

func (servicePool *ServicePool) UpdateService(service *Service) {

	servicePool.mux.Lock()
	defer servicePool.mux.Unlock()

	known, index := servicePool.FirstKnownInstance(service)

	if known {
		servicePool.Services[index].WasSeen(time.Now())
	} else {
		servicePool.Services = append(servicePool.Services, service)
	}
}

func (servicePool *ServicePool) FirstKnownInstance(service *Service) (bool, int) {

	b := service.URL
	for i, member := range servicePool.Services {
		a := member.URL
		if *a == *b {
			return true, i
		}
	}
	return false, -1
}
