package main

import (
	"time"

	"github.com/Festivals-App/festivals-gateway/server"
	"github.com/Festivals-App/festivals-gateway/server/config"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Info().Msg("Server startup")

	root := servertools.ContainerPathArgument()
	configFilePath := root + "/etc/festivals-gateway.conf"

	conf := config.ParseConfig(configFilePath)
	log.Info().Msg("Server configuration was initialized")

	servertools.InitializeGlobalLogger(conf.InfoLog, true)
	log.Info().Msg("Logger initialized")

	server := server.NewServer(conf)
	go server.Run(conf)
	log.Info().Msg("Server did start")

	go sendHeartbeat(conf)
	log.Info().Msg("Heartbeat routine was started")

	// wait forever
	// https://stackoverflow.com/questions/36419054/go-projects-main-goroutine-sleep-forever
	select {}
}

func sendHeartbeat(conf *config.Config) {

	heartbeatClient, err := servertools.HeartbeatClient(conf.TLSCert, conf.TLSKey, conf.TLSRootCert)
	if err != nil {
		log.Fatal().Err(err).Str("type", "server").Msg("Failed to create heartbeat client")
	}
	var beat *servertools.Heartbeat = &servertools.Heartbeat{Service: "festivals-gateway", Host: "https://gateway." + conf.ServiceBindHost, Port: conf.ServicePort, Available: true}

	t := time.NewTicker(time.Duration(conf.Interval) * time.Second)
	defer t.Stop()
	for range t.C {
		err = servertools.SendHeartbeat(heartbeatClient, conf.LoversEar, conf.ServiceKey, beat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send heartbeat")
		}
	}
}
