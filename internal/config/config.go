package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds the configuration values
type Config struct {
	LogLevel       string
	BrokerAddress  string
	BrokerPort     int
	BrokerPassword string
}

// Options holds the configuration options for building a Config.
type Options struct {
	logLevel       string
	brokerAddress  string
	brokerPort     int
	brokerPassword string
}

// NewConfig creates a new Config struct and applies the provided options.
func NewConfig(opts *Options) *Config {
	return &Config{
		LogLevel:       opts.logLevel,
		BrokerAddress:  opts.brokerAddress,
		BrokerPort:     opts.brokerPort,
		BrokerPassword: opts.brokerPassword,
	}
}

// NewConfigOpt initializes a new Options struct with default values.
func NewConfigOpt() *Options {
	return &Options{}
}

// WithLogLevel sets the log level in the Options.
func (o *Options) WithLogLevel() *Options {
	o.logLevel = getEnv("LOG_LEVEL")
	return o
}

// WithBrokerAddress sets the broker address in the Options
func (o *Options) WithBrokerAddress() *Options {
	o.brokerAddress = getEnv("REDIS_ADDR")
	return o
}

// WithBrokerPort sets the broker port in the Options
func (o *Options) WithBrokerPort() *Options {
	o.brokerPort, _ = strconv.Atoi(getEnv("REDIS_PORT"))
	return o
}

// WithBrokerPassword sets the broker password in the Options
func (o *Options) WithBrokerPassword() *Options {
	o.brokerPassword = getEnv("REDIS_PASSWORD")
	return o
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Panicf("Environment variable %s is not set", key)
	}

	return value
}
