package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerHostPort string `env:"SM_SERVER_HOST_PORT"`
	DSN            string `env:"SM_DATA_SOURCE_NAME"`
}

func newConfig() *Config {
	return &Config{
		ServerHostPort: ServerHostPort,
		DSN:            DataSourceName,
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
	args := os.Args
	fmt.Printf("All arguments: %v\n", args)
	var port string
	flag.StringVar(&port, "p", c.ServerHostPort, "host to listen on")
	c.ServerHostPort = ":" + port
	flag.StringVar(&c.DSN, "d", c.DSN, "data storage name")
	flag.Parse()
}
