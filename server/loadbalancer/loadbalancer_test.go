package loadbalancer_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/loadbalancer"
)

func TestFirstKnownInstance(t *testing.T) {

	servicePool := testPool()
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

func testPool() *loadbalancer.ServicePool {

	servicePool := loadbalancer.ServicePool{
		Services: []*loadbalancer.Service{gatewayService1(), festivalsService1()},
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
