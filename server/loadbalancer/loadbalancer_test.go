package loadbalancer_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/loadbalancer"
)

func TestFirstKnownInstance(t *testing.T) {

	servicePool := testPool()
	servicePool.Services = append(servicePool.Services, gatewayService1())
	servicePool.Services = append(servicePool.Services, festivalsService1())

	known, index := servicePool.FirstKnownInstance(gatewayService2())
	if known == false {
		t.Errorf("Coulden't find service instance in service pool, but expected to find one service with index 0.")
	}
	if index != 0 {
		t.Errorf("Did find service instance in service pool, but expected to find service with index 0 not %d.", index)
	}

	known, index = servicePool.FirstKnownInstance(festivalsfileService1())
	if known == true {
		t.Errorf("Did find service instance in service pool, but expected to not find one.")
	}
	if index != -1 {
		t.Errorf("Did not find service instance in service pool, but expected index to be -1 not %d.", index)
	}
}

func IsEqualToService(t *testing.T) {

	festival1 := festivalsService1()
	festival2 := festivalsService2()

	if festival1.IsEqualTo(festival2) {
		t.Errorf("Service (%v) and service (%v) should not be equal.", festival1, festival2)
	}

	if !festival1.IsEqualTo(festival1) {
		t.Errorf("Service (%v) and service (%v) should be equal.", festival1, festival1)
	}
}

func TestGetNextService(t *testing.T) {

	servicePool := testPool()

	festival1 := festivalsService1()
	festival2 := festivalsService2()
	festival3 := festivalsService3()
	festival4_stale := festivalsService4_stale()
	servicePool.Services = append(servicePool.Services, festival1)

	service := servicePool.GetNextService()
	if service != festival1 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival1, service)
	}
	service = servicePool.GetNextService()
	if service != festival1 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival1, service)
	}
	service = servicePool.GetNextService()
	if service != festival1 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival1, service)
	}

	servicePool = testPool()
	servicePool.Services = append(servicePool.Services, festival1)
	servicePool.Services = append(servicePool.Services, festival2)

	service = servicePool.GetNextService()
	if service != festival2 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival2.URL, service.URL)
	}
	service = servicePool.GetNextService()
	if service != festival1 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival1.URL, service.URL)
	}
	service = servicePool.GetNextService()
	if service != festival2 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival2.URL, service.URL)
	}
	service = servicePool.GetNextService()
	if service != festival1 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival1.URL, service.URL)
	}
	service = servicePool.GetNextService()
	if service != festival2 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival2.URL, service.URL)
	}

	servicePool = testPool()
	servicePool.Services = append(servicePool.Services, festival1)
	servicePool.Services = append(servicePool.Services, festival2)
	servicePool.Services = append(servicePool.Services, festival3)

	service = servicePool.GetNextService()
	if service != festival2 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival2.URL, service.URL)
	}

	service = servicePool.GetNextService()
	if service != festival3 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival3.URL, service.URL)
	}

	service = servicePool.GetNextService()
	if service != festival1 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival1.URL, service.URL)
	}

	servicePool = testPool()
	servicePool.Services = append(servicePool.Services, festival1)
	servicePool.Services = append(servicePool.Services, festival2)
	servicePool.Services = append(servicePool.Services, festival3)
	servicePool.Services = append(servicePool.Services, festival4_stale)

	servicePool.Current = 3

	service = servicePool.GetNextService()
	if service != festival1 {
		t.Errorf("Next service should be service (%v) but is service (%v)", festival1.URL, service.URL)
	}
}

func testPool() *loadbalancer.ServicePool {

	servicePool := loadbalancer.ServicePool{
		Services: []*loadbalancer.Service{},
		Current:  0,
	}
	return &servicePool
}

func gatewayService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://gateway.festivalsapp.home:443")
	return &loadbalancer.Service{
		Name:  "festivals-gateway",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func gatewayService2() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://gateway.festivalsapp.home:443")
	return &loadbalancer.Service{
		Name:  "festivals-gateway",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func festivalsService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://festivals-0.festivalsapp.home:10439")
	return &loadbalancer.Service{
		Name:  "festivals-server",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func festivalsService2() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://festivals-1.festivalsapp.home:10439")
	return &loadbalancer.Service{
		Name:  "festivals-server",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func festivalsService3() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://festivals-2.festivalsapp.home:10439")
	return &loadbalancer.Service{
		Name:  "festivals-server",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func festivalsService4_stale() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://festivals-3.festivalsapp.home:10439")
	return &loadbalancer.Service{
		Name:  "festivals-server",
		URL:   serviceurl,
		At:    time.Now().Add(time.Second * 40),
		Alive: true,
	}
}

func festivalsService5_stale() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://festivals-4.festivalsapp.home:10439")
	return &loadbalancer.Service{
		Name:  "festivals-server",
		URL:   serviceurl,
		At:    time.Now().Add(time.Second * 40),
		Alive: true,
	}
}

func festivalsfileService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://fileserver-0.festivalsapp.home:1910")
	return &loadbalancer.Service{
		Name:  "festivals-fileserver",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func identityService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://identity-0.festivalsapp.home:22580")
	return &loadbalancer.Service{
		Name:  "festivals-identity-server",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func databaseNodeService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://database-0.festivalsapp.home:22397")
	return &loadbalancer.Service{
		Name:  "festivals-database-node",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func databaseService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("mysql://192.168.8.117:3306")
	return &loadbalancer.Service{
		Name:  "festivals-database",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func websiteNodeService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("https://website-0.festivalsapp.home:48155")
	return &loadbalancer.Service{
		Name:  "festivals-website-node",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}

func websiteService1() *loadbalancer.Service {
	serviceurl, _ := url.Parse("http://website-0.festivalsapp.home:8080")
	return &loadbalancer.Service{
		Name:  "festivals-website",
		URL:   serviceurl,
		At:    time.Now(),
		Alive: true,
	}
}
