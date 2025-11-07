package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Addr            string `env:"ADDRESS"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
}

func InitServerConfig() (*ServerConfig, error) {

	serverConfig := new(ServerConfig)

	flagAddr := flag.String("a", "localhost:8080", "Metrics server address")
	flagDB := flag.String("d", "", "Database DSN")
	flagFileStoragePath := flag.String("f", "./tmp/metrics.json", "Metrics store filepath")
	flagRestore := flag.Bool("r", false, "Flag to load previous values from file on startup")
	flagStoreInterval := flag.Int("i", 300, "Metrics store interval in seconds (0 to sync)")
	flag.Parse()

	serverConfig.Addr = *flagAddr
	serverConfig.DatabaseDSN = *flagDB
	serverConfig.FileStoragePath = *flagFileStoragePath
	serverConfig.Restore = *flagRestore
	serverConfig.StoreInterval = *flagStoreInterval

	err := env.Parse(serverConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot parse env variables: %w", err)
	}

	if serverConfig.StoreInterval < 0 {
		return nil, fmt.Errorf("bad store interval (<0): %w", errors.New("metrics store interval must be not less than zero"))
	}

	return serverConfig, nil
}
