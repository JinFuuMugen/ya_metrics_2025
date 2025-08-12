package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/logger"
	models "github.com/JinFuuMugen/ya_metrics_2025/internal/model"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Errorf("method %s is not allowed", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metricName := r.PathValue(`metric_name`)
	if metricName == "" {
		logger.Errorf("empty metric name %s", metricName)
		http.Error(w, "no metric name provided", http.StatusNotFound)
		return
	}

	metricType := r.PathValue(`metric_type`)
	metricValue := r.PathValue(`metric_value`)

	switch metricType {

	case models.Gauge:

		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			logger.Errorf("cannot parse %s to float: %s", metricValue, err)
			http.Error(w, fmt.Sprintf("invalid value: %s of type %s", metricValue, metricType), http.StatusBadRequest)
			return
		}

		storage.SetGauge(metricName, floatValue)

	case models.Counter:

		intValue, err := strconv.Atoi(metricValue)
		if err != nil {
			logger.Errorf("cannot parse %s to int: %s", metricValue, err)
			http.Error(w, fmt.Sprintf("invalid value: %s of type %s", metricValue, metricType), http.StatusBadRequest)
			return
		}

		storage.AddCounter(metricName, int64(intValue))

	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		logger.Errorf("unknown metric type: %s", metricType)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		logger.Errorf("method %s is not allowed", r.Method)
	}

	metricType := r.PathValue(`metric_type`)
	metricName := r.PathValue(`metric_name`)

	switch metricType {
	case storage.MetricTypeCounter:
		counter, err := storage.GetCounter(metricName)
		if err != nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			logger.Errorf("metric not found: %s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(counter.GetValueString()))

	case storage.MetricTypeGauge:
		gauge, err := storage.GetGauge(metricName)
		if err != nil {
			http.Error(w, "metric not found", http.StatusNotFound)
			logger.Errorf("metric not found: %s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(gauge.GetValueString()))

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		logger.Errorf("unknown metric type: %s", metricType)
		return
	}

}

func UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		logger.Errorf("method %s is not allowed", r.Method)
		return
	}

	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		logger.Errorf("cannot read body: %s", err)
		return
	}

	var metric models.Metrics

	if err := json.Unmarshal(buffer.Bytes(), &metric); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		logger.Errorf("cannot unmarshal body: %s", err)
	}

	switch metric.MType {

	case storage.MetricTypeCounter:
		if metric.ID != "" && metric.Delta != nil {
			storage.AddCounter(metric.ID, *metric.Delta)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			logger.Errorf("missing crucial fields in request %s", buffer.String())
			return
		}

	case storage.MetricTypeGauge:
		if metric.ID != "" && metric.Value != nil {
			storage.SetGauge(metric.ID, *metric.Value)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			logger.Errorf("missing crucial fields in request %s", buffer.String())
			return
		}

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		logger.Errorf("unknown metric type: %s", metric.MType)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		logger.Errorf("method %s is not allowed", r.Method)
		return
	}

	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		logger.Errorf("cannot read body: %s", err)
		return
	}

	var metric models.Metrics

	if err := json.Unmarshal(buffer.Bytes(), &metric); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		logger.Errorf("cannot unmarshal body: %s", err)
	}

	switch metric.MType {
	case storage.MetricTypeCounter:
		if metric.ID != "" {
			c, err := storage.GetCounter(metric.ID)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				logger.Errorf("metric not found: %s", err)
				return
			}

			metric.Delta = &c.Value

		} else {
			http.Error(w, "empty metric name", http.StatusBadRequest)
			logger.Errorf("empty metric name %s", metric.ID)
			return
		}

	case storage.MetricTypeGauge:
		if metric.ID != "" {
			g, err := storage.GetGauge(metric.ID)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				logger.Errorf("metric not found: %s", err)
				return
			}

			metric.Value = &g.Value

		} else {
			http.Error(w, "empty metric name", http.StatusBadRequest)
			logger.Errorf("empty metric name %s", metric.ID)
			return
		}

	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		logger.Errorf("unknown metric type: %s", metric.MType)
		return
	}

	marshMetric, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		logger.Errorf("cannot marshal metric: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.Write(marshMetric)

}
