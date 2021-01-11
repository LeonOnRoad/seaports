package config

import (
	"time"

	"github.com/spf13/viper"
)

const envPort = "PORT"
const envShutdownTimeout = "SHUTDOWN_TIMEOUT"
const envRedisEndpoint = "REDIS_ENDPOINT"

type Config struct {
	Port            int
	ShutdownTimeout time.Duration
	RedisEndpoint   string
}

// LoadConfig already includes default values in case env-vars are not provided
func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.SetDefault(envPort, 8081)
	viper.SetDefault(envShutdownTimeout, "10s")
	viper.SetDefault(envRedisEndpoint, "localhost:6379") // if this is removed, then in memory repository will be used

	config := &Config{}
	config.Port = viper.GetInt(envPort)
	config.RedisEndpoint = viper.GetString(envRedisEndpoint)

	var err error
	config.ShutdownTimeout, err = time.ParseDuration(viper.GetString(envShutdownTimeout))
	if err != nil {
		config.ShutdownTimeout = 10 * time.Second
	}
	return config
}
