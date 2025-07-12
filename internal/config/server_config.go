package config

import "flag"

type ServerConfig struct {
	Addr string
}

func InitServerConfig() *ServerConfig {

	addr := flag.String("a", "localhost:8080", "Metrics server address")

	flag.Parse()

	return &ServerConfig{Addr: *addr}
}
