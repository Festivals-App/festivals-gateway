package config

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	ServiceBindAddress string
	ServiceBindHost    string
	ServicePort        int
	APIKeys            []string
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
		log.Fatal("server initialize: could not read default config file.")
	}
	path = path + "/config_template.toml"
	return ParseConfig(path)
}

func ParseConfig(cfgFile string) *Config {

	content, err := toml.LoadFile(cfgFile)
	if err != nil {
		log.Fatal("server initialize: could not read config file at '" + cfgFile + "'. Error: " + err.Error())
	}

	serverBindAdress := content.Get("service.bind-address").(string)
	serverBindHost := content.Get("service.bind-host").(string)
	serverPort := content.Get("service.port").(int64)

	keyValues := content.Get("authentication.api-keys").([]interface{})
	keys := make([]string, len(keyValues))
	for i, v := range keyValues {
		keys[i] = v.(string)
	}

	return &Config{
		ServiceBindAddress: serverBindAdress,
		ServiceBindHost:    serverBindHost,
		ServicePort:        int(serverPort),
		APIKeys:            keys,
	}
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