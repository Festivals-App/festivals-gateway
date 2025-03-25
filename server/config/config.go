package config

import (
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"

	"github.com/pelletier/go-toml"
)

type Config struct {
	ServiceBindHost  string
	ServicePort      int
	ServiceKey       string
	TLSRootCert      string
	TLSCert          string
	TLSKey           string
	LoversEar        string
	Interval         int
	IdentityEndpoint string
	InfoLog          string
	TraceLog         string
}

func ParseConfig(cfgFile string) *Config {

	content, err := toml.LoadFile(cfgFile)
	if err != nil {
		log.Fatal().Msg("server initialize: could not read config file at '" + cfgFile + "'. Error: " + err.Error())
	}

	serviceBindHost := content.Get("service.bind-host").(string)
	servicePort := content.Get("service.port").(int64)
	serviceKey := content.Get("service.key").(string)

	tlsrootcert := content.Get("tls.festivaslapp-root-ca").(string)
	tlscert := content.Get("tls.cert").(string)
	tlskey := content.Get("tls.key").(string)

	loversear := content.Get("heartbeat.endpoint").(string)
	interval := content.Get("heartbeat.interval").(int64)

	identity := content.Get("authentication.endpoint").(string)

	infoLogPath := content.Get("log.info").(string)
	traceLogPath := content.Get("log.trace").(string)

	tlsrootcert = servertools.ExpandTilde(tlsrootcert)
	tlscert = servertools.ExpandTilde(tlscert)
	tlskey = servertools.ExpandTilde(tlskey)
	infoLogPath = servertools.ExpandTilde(infoLogPath)
	traceLogPath = servertools.ExpandTilde(traceLogPath)

	return &Config{
		ServiceBindHost:  serviceBindHost,
		ServicePort:      int(servicePort),
		ServiceKey:       serviceKey,
		TLSRootCert:      tlsrootcert,
		TLSCert:          tlscert,
		TLSKey:           tlskey,
		LoversEar:        loversear,
		Interval:         int(interval),
		IdentityEndpoint: identity,
		InfoLog:          infoLogPath,
		TraceLog:         traceLogPath,
	}
}
