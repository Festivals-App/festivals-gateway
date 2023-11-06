package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Festivals-App/festivals-gateway/server"
	"github.com/Festivals-App/festivals-gateway/server/config"
	"github.com/Festivals-App/festivals-gateway/server/heartbeat"
	"github.com/Festivals-App/festivals-gateway/server/logger"
	"github.com/rs/zerolog/log"
)

func main() {

	logger.InitializeGlobalLogger("/var/log/festivals-gateway/info.log", true)

	log.Info().Msg("Server startup.")

	conf := config.DefaultConfig()

	log.Info().Msg("Server configuration was initialized.")

	config.CheckForArguments()

	serverInstance := &server.Server{}
	serverInstance.Initialize(conf)

	// Redirect traffic from port 80 to used port
	go http.ListenAndServe(":80", serverInstance.CertManager.HTTPHandler(nil))
	log.Info().Msg("Start redirecting http traffic from port 80 to port " + strconv.Itoa(serverInstance.Config.ServicePort) + " and https")

	go serverInstance.Run(conf)
	log.Info().Msg("Server did start.")

	go sendHeartbeat(conf)
	log.Info().Msg("Heartbeat routine was started.")

	// wait forever
	// https://stackoverflow.com/questions/36419054/go-projects-main-goroutine-sleep-forever
	select {}
}

func sendHeartbeat(conf *config.Config) {
	for {
		timer := time.After(time.Second * time.Duration(conf.Interval))
		<-timer
		var beat *heartbeat.Heartbeat = &heartbeat.Heartbeat{Service: "festivals-gateway", Host: "https://gateway." + conf.ServiceBindHost, Port: conf.ServicePort, Available: true}
		err := heartbeat.SendHeartbeat(conf.LoversEar, conf.ServiceKey, conf.TLSCert, conf.TLSKey, conf.TLSRootCert, beat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send heartbeat")
		}
	}
}
