package api

import (
	"github.com/JinFuuMugen/ya_metrics_2025/internal/handler"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/logger"
	customMiddleware "github.com/JinFuuMugen/ya_metrics_2025/internal/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRouter() *chi.Mux {
	rout := chi.NewRouter()

	rout.Use(logger.LoggerMiddleware)
	rout.Use(middleware.StripSlashes)
	rout.Use(customMiddleware.GzipMiddleware)

	rout.Post("/update/{metric_type}/{metric_name}/{metric_value}", handler.UpdateMetricHandler)
	rout.Post("/update", handler.UpdateMetricJSONHandler)
	rout.Get("/value/{metric_type}/{metric_name}", handler.GetMetricHandler)
	rout.Post("/value", handler.GetMetricJSONHandler)
	rout.Get("/", handler.InfoPageHandler)

	return rout
}
