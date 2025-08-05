package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func InitServerConfig() (*ServerConfig, error) {

	serverConfig := new(ServerConfig)

	flagAddr := flag.String("a", "localhost:8080", "Metrics server address")
	flagStoreInterval := flag.Int("i", 300, "Metrics store interval in seconds (0 to sync)")
	flagFileStoragePath := flag.String("f", "/data/metrics.txt", "Metrics store filepath")
	flagRestore := flag.Bool("r", false, "Flag to load previous values from file on startup")
	flag.Parse()

	serverConfig.Addr = *flagAddr
	serverConfig.StoreInterval = *flagStoreInterval
	serverConfig.FileStoragePath = *flagFileStoragePath
	serverConfig.Restore = *flagRestore

	err := env.Parse(serverConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot parse env variables: %w", err)
	}

	return serverConfig, nil
}
