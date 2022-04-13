package main

import (
	"net"
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

	logger.Initialize("/var/log/festivals-gateway/info.log", true)

	log.Info().Msg("Server startup.")

	conf := config.DefaultConfig()
	//if len(os.Args) > 1 {
	//	conf = config.ParseConfig(os.Args[1])
	//}

	log.Info().Msg("Server configuration was initialized.")

	config.CheckForArguments()

	serverInstance := &server.Server{}
	serverInstance.Initialize(conf)

	// Redirect traffic from port 80 to used port
	go http.ListenAndServe(":80", serverInstance.CertManager.HTTPHandler(nil))
	log.Info().Msg("Start redirecting http traffic from port 80 to port " + strconv.Itoa(serverInstance.Config.ServicePort) + " and https")

	go serverInstance.Run(":" + strconv.Itoa(conf.ServicePort))
	log.Info().Msg("Server did start.")

	go sendHeartbeat(conf)
	log.Info().Msg("Heartbeat routine was started.")

	// wait forever
	// https://stackoverflow.com/questions/36419054/go-projects-main-goroutine-sleep-forever
	select {}
}

func sendHeartbeat(conf *config.Config) {
	for {
		timer := time.After(time.Second * 2)
		<-timer
		var beat *heartbeat.Heartbeat = &heartbeat.Heartbeat{Service: "festivals-gateway", Host: conf.ServiceBindHost, Port: conf.ServicePort, Available: true}
		err := heartbeat.SendHeartbeat(conf.LoversEar, conf.ServiceKey, beat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send heartbeat")
		}
	}
}

func redirectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Please use HTTPS for your request", http.StatusBadRequest)
		return
	}
	target := "https://" + stripPort(r.Host) + r.URL.RequestURI()
	log.Info().Msg("Redirecting request for '" + r.URL.String() + "' to HTTPS")
	http.Redirect(w, r, target, http.StatusFound)
}

func stripPort(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}
	return net.JoinHostPort(host, "443")
}
