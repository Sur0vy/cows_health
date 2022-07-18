package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHostPort string `env:"SM_SERVER_HOST_PORT"`
	DSN            string `env:"SM_DATA_SOURCE_NAME"`
	IsDebug        bool   `env:"IS_DEBUG"`
	LogFile        string `env:"LOG_FILE"`
}

func newConfig() *Config {
	return &Config{
		ServerHostPort: ServerHostPort,
		DSN:            DataSourceName,
		IsDebug:        false,
	}
}

type option func(*Config)

func LoadParams() option {
	return func(c *Config) {
		c.flagParams()
		_ = env.Parse(c)
	}
}

func Setup(opts ...option) *Config {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Config) flagParams() {
	var port string
	flag.StringVar(&port, "p", c.ServerHostPort, "host to listen on")
	flag.StringVar(&c.DSN, "d", c.DSN, "data storage name")
	flag.BoolVar(&c.IsDebug, "deb", c.IsDebug, "debug logging mode")
	flag.StringVar(&c.LogFile, "lg", "", "log file")
	flag.Parse()
	c.ServerHostPort = ":" + port
}
