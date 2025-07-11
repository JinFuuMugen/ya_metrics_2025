package config

import "flag"

type ServerConfig struct {
	Addr string
}

func InitConfig() *ServerConfig {

	addr := flag.String("Server address", "localhost:8080", "Metrics server address")

	return &ServerConfig{Addr: *addr}
}
