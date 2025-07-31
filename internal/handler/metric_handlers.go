package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/JinFuuMugen/ya_metrics_2025/internal/model"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metricName := r.PathValue(`metric_name`)
	if metricName == "" {
		http.Error(w, "no metric name provided", http.StatusNotFound)
		return
	}

	metricType := r.PathValue(`metric_type`)
	metricValue := r.PathValue(`metric_value`)

	switch metricType {

	case models.Gauge:

		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid value: %s of type %s", metricValue, metricType), http.StatusBadRequest)
			return
		}

		storage.SetGauge(metricName, floatValue)

	case models.Counter:

		intValue, err := strconv.Atoi(metricValue)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid value: %s of type %s", metricValue, metricType), http.StatusBadRequest)
			return
		}

		storage.AddCounter(metricName, int64(intValue))

	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	metricType := r.PathValue(`metric_type`)
	metricName := r.PathValue(`metric_name`)

	switch metricType {
	case storage.MetricTypeCounter:
		counter, err := storage.GetCounter(metricName)
		if err != nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(counter.GetValueString()))

	case storage.MetricTypeGauge:
		gauge, err := storage.GetGauge(metricName)
		if err != nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(gauge.GetValueString()))

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

}

func UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var metric models.Metrics

	if err := json.Unmarshal(buffer.Bytes(), &metric); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	switch metric.MType {

	case storage.MetricTypeCounter:
		if metric.ID != "" && metric.Delta != nil {
			storage.AddCounter(metric.ID, *metric.Delta)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

	case storage.MetricTypeGauge:
		if metric.ID != "" && metric.Value != nil {
			storage.SetGauge(metric.ID, *metric.Value)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
