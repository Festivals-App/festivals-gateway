package main

import (
	"os"
	"strconv"

	"github.com/Festivals-App/festivals-gateway/server"
	"github.com/Festivals-App/festivals-gateway/server/config"
)

func main() {
	conf := config.DefaultConfig()
	if len(os.Args) > 1 {
		conf = config.ParseConfig(os.Args[1])
	}

	serverInstance := &server.Server{}
	serverInstance.Initialize(conf)
	serverInstance.Run(conf.ServiceBindAddress + ":" + strconv.Itoa(conf.ServicePort))
}
