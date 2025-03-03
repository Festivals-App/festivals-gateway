package config

import (
	"os"

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
}

func DefaultConfig() *Config {

	// first we try to parse the config at the global configuration path
	if fileExists("/etc/festivals-gateway.conf") {
		config := ParseConfig("/etc/festivals-gateway.conf")
		if config != nil {
			return config
		}
	}

	// if there is no global configuration check the current folder for the template config file
	// this is mostly so the application will run in development environment
	path, err := os.Getwd()
	if err != nil {
		log.Fatal().Msg("server initialize: could not read default config file.")
	}
	path = path + "/config_template.toml"
	return ParseConfig(path)
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

	checkForDebugMode()

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
	}
}

func checkForDebugMode() {

	if len(os.Args) == 2 {
		if os.Args[1] == "--debug" {

			os.Setenv("DEBUG", "true")
			log.Info().Msg("Running in debug mode")
		}
	}
}

func IsRunningInDebug() bool {
	_, isPresent := os.LookupEnv("DEBUG")
	return isPresent
}

func IsRunningInProduction() bool {
	_, isPresent := os.LookupEnv("DEBUG")
	return !isPresent
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
// see: https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
