package config

import (
	"github.com/rs/zerolog/log"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	ServiceBindAddress string
	ServiceBindHost    string
	ServicePort        int
	ServiceKey         string
	Website            string
	TLSCert            string
	TLSKey             string
	LoversEar          string
	APIKeys            []string
	AdminKeys          []string
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

	serviceBindAdress := content.Get("service.bind-address").(string)
	serviceBindHost := content.Get("service.bind-host").(string)
	servicePort := content.Get("service.port").(int64)
	serviceKey := content.Get("service.key").(string)
	website := content.Get("service.website").(string)

	tlscert := content.Get("tls.cert").(string)
	tlskey := content.Get("tls.key").(string)

	loversear := content.Get("heartbeat.endpoint").(string)

	keyValues := content.Get("authentication.api-keys").([]interface{})
	keys := make([]string, len(keyValues))
	for i, v := range keyValues {
		keys[i] = v.(string)
	}
	adminKeyValues := content.Get("authentication.admin-keys").([]interface{})
	adminKeys := make([]string, len(adminKeyValues))
	for i, v := range adminKeyValues {
		adminKeys[i] = v.(string)
	}

	return &Config{
		ServiceBindAddress: serviceBindAdress,
		ServiceBindHost:    serviceBindHost,
		ServicePort:        int(servicePort),
		ServiceKey:         serviceKey,
		TLSCert:            tlscert,
		TLSKey:             tlskey,
		Website:            website,
		LoversEar:          loversear,
		APIKeys:            keys,
		AdminKeys:          adminKeys,
	}
}

func CheckForArguments() {

	if len(os.Args) == 2 {
		if os.Args[1] == "--debug" {

			os.Setenv("DEBUG", "true")
			log.Info().Msg("Running in debug mode")
		}
	}
}

func Debug() bool {
	_, isPresent := os.LookupEnv("DEBUG")
	return isPresent
}

func Production() bool {
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
