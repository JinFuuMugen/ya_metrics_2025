package api

import (
	"github.com/JinFuuMugen/ya_metrics_2025/internal/handler"
	"github.com/go-chi/chi/v5"
)

func InitRouter() *chi.Mux {
	rout := chi.NewRouter()

	rout.Post("/update/{metric_type}/{metric_name}/{metric_value}", handler.UpdateMetricHandler)
	rout.Get("/value/{metric_type}/{metric_name}", handler.GetMetricHandler)
	rout.Get("/", handler.InfoPageHandler)

	return rout
}
