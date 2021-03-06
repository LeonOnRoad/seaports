package config

import (
	"time"

	"github.com/spf13/viper"
)

const envPort = "PORT"
const envPortsServiceEndpoint = "PORTS_SERVICE_ENDPOINT"
const envShutdownTimeout = "SHUTDOWN_TIMEOUT"

type Config struct {
	Port                 int
	PortsServiceEndpoint string
	ShutdownTimeout      time.Duration
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.SetDefault(envPort, 8080)
	viper.SetDefault(envPortsServiceEndpoint, "localhost:8081")
	viper.SetDefault(envShutdownTimeout, "10s")

	config := &Config{}
	config.Port = viper.GetInt(envPort)
	config.PortsServiceEndpoint = viper.GetString(envPortsServiceEndpoint)
	var err error
	config.ShutdownTimeout, err = time.ParseDuration(viper.GetString(envShutdownTimeout))
	if err != nil {
		config.ShutdownTimeout = 10 * time.Second
	}
	return config
}
