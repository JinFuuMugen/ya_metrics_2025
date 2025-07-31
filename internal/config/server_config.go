package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Addr string `env:"ADDRESS"`
}

func InitServerConfig() (*ServerConfig, error) {

	serverConfig := new(ServerConfig)

	flagAddr := flag.String("a", "localhost:8080", "Metrics server address")
	flag.Parse()

	serverConfig.Addr = *flagAddr

	err := env.Parse(serverConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot parse env variables: %w", err)
	}

	return serverConfig, nil
}
