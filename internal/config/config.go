package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHostPort string `env:"SM_SERVER_HOST_PORT"`
	FrontDir       string `env:"SM_FRONT_DIR"`
	DSN            string `env:"SM_DATA_SOURCE_NAME"`
	//CurrentUser     string
	//CurrentUserHash string
}

func newConfig() *Config {
	return &Config{
		ServerHostPort: ServerHostPort,
		FrontDir:       FrontendDir,
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
	flag.StringVar(&c.FrontDir, "f", c.FrontDir, "frontend setup directory")
	flag.StringVar(&c.DSN, "d", c.DSN, "data storage name")
	flag.Parse()
}
