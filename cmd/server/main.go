package main

import (
	"fmt"
	"net/http"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/handler"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

func main() {

	storage.InitStorage()

	mux := http.NewServeMux()

	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", handler.UpdateMetricHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(fmt.Errorf("cannot start server: %w", err))
	}
}
