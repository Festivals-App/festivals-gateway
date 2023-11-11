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

func (servicePool *ServicePool) GetNextService() *Service {

	servicePool.mux.Lock()
	defer servicePool.mux.Unlock()

	// loop entire backends to find out an Alive backend
	next := servicePool.NextIndex()

	numberOfServices := len(servicePool.Services)
	l := numberOfServices + next // start from next and move a full cycle

	for i := next; i < l; i++ {
		idx := i % numberOfServices // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one
		if servicePool.Services[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&servicePool.Current, uint64(idx)) // mark the current one
			}
			return servicePool.Services[idx]
		}
	}
	return nil
}

func (servicePool *ServicePool) NextIndex() int {
	return int(atomic.AddUint64(&servicePool.Current, uint64(1)) % uint64(len(servicePool.Services)))
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
