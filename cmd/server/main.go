package main

import (
	"fmt"
	"net/http"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/api"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/config"
)

func main() {

	cfg, err := config.InitServerConfig()
	if err != nil {
		panic(fmt.Errorf("cannot init server config: %w", err))
	}

	rout := api.InitRouter()

	err = http.ListenAndServe(cfg.Addr, rout)
	if err != nil {
		panic(fmt.Errorf("cannot start server: %w", err))
	}
}
