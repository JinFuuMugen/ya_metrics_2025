package main

import (
	"context"
	defaultLogger "log"
	"net/http"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/api"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/config"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/filestorage"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/logger"
)

func main() {

	err := logger.InitLogger()
	if err != nil {
		defaultLogger.Fatalf("cannot init custom logger: %s", err)
	}
	defer logger.Sync()

	cfg, err := config.InitServerConfig()
	if err != nil {
		logger.Fatalf("cannot init server config: %w", err)
	}

	rout := api.InitRouter(cfg)

	arbitrator, err := filestorage.InitArbitrator(cfg)
	if err != nil {
		logger.Fatalf("cannot init arbitrator: %w", err)
	}

	err = arbitrator.StartArbitrator(context.Background())
	if err != nil {
		logger.Fatalf("cannot start arbitrator: %w", err)
	}

	err = http.ListenAndServe(cfg.Addr, rout)
	if err != nil {
		logger.Fatalf("cannot start server: %w", err)
	}
}
