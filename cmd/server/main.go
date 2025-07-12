package main

import (
	"fmt"
	"net/http"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/config"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.InitServerConfig()

	rout := chi.NewRouter()

	rout.Post("/update/{metric_type}/{metric_name}/{metric_value}", handler.UpdateMetricHandler)
	rout.Get("/value/{metric_type}/{metric_name}", handler.GetMetricHandler)
	rout.Get("/", handler.InfoPageHandler)

	err := http.ListenAndServe(cfg.Addr, rout)
	if err != nil {
		panic(fmt.Errorf("cannot start server: %w", err))
	}
}
