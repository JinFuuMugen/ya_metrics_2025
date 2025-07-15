package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

func TestUpdateMetricHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		pathVals    map[string]string
		wantStatus  int
		wantGauge   float64
		wantCounter int64
	}{

		{
			name:       "valid gauge",
			method:     http.MethodPost,
			url:        "/update/gauge/GaugeMetr/123.4",
			pathVals:   map[string]string{"metric_name": "GaugeMetr", "metric_type": "gauge", "metric_value": "123.4"},
			wantStatus: http.StatusOK,
			wantGauge:  123.4,
		},
		{
			name:       "wrong method",
			method:     http.MethodGet,
			url:        "/update/gauge/GaugeMetr/123.4",
			pathVals:   map[string]string{"metric_name": "GaugeMetr", "metric_type": "gauge", "metric_value": "123.4"},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "wrong value",
			method:     http.MethodPost,
			url:        "/update/gauge/GaugeMetr/1e3",
			pathVals:   map[string]string{"metric_name": "GaugeMetr", "metric_type": "gauge", "metric_value": "1e3a"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "valid counter",
			method:      http.MethodPost,
			url:         "/update/counter/CounterMetr/2",
			pathVals:    map[string]string{"metric_name": "CounterMetr", "metric_type": "counter", "metric_value": "2"},
			wantStatus:  http.StatusOK,
			wantCounter: 2,
		},
		{
			name:       "float counter",
			method:     http.MethodPost,
			pathVals:   map[string]string{"metric_name": "CounterMetr", "metric_type": "counter", "metric_value": "123.4"},
			url:        "/update/counter/CounterMetr/123.4",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(tt.method, tt.url, nil)
			for k, v := range tt.pathVals {
				r.SetPathValue(k, v)
			}

			w := httptest.NewRecorder()

			UpdateMetricHandler(w, r)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
				return
			}

			if tt.wantGauge != 0 {
				got, err := storage.GetGauge("GaugeMetr")
				if err != nil {
					t.Errorf("error getting gauge value: %s", err)
					return
				}

				if got.GetValue() != tt.wantGauge {
					t.Errorf("gauge = %v, want %v", got.GetValue(), tt.wantGauge)
					return
				}
			}

			if tt.wantCounter != 0 {
				got, err := storage.GetCounter("CounterMetr")
				if err != nil {
					t.Errorf("error getting counter value: %s", err)
					return
				}

				if got.GetValue() != tt.wantCounter {
					t.Errorf("counter = %v, want %v", got, tt.wantGauge)
					return
				}
			}
		})
	}
}

func TestGetMetricHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		pathVals    map[string]string
		wantStatus  int
		wantGauge   float64
		wantCounter int64
	}{
		{
			name:       "valid gauge",
			method:     http.MethodGet,
			url:        "/value/gauge/valid_gauge",
			pathVals:   map[string]string{"metric_name": "valid_gauge", "metric_type": "gauge"},
			wantGauge:  100.100,
			wantStatus: http.StatusOK,
		},
		{
			name:        "valid counter",
			method:      http.MethodGet,
			url:         "/value/counter/valid_counter",
			pathVals:    map[string]string{"metric_name": "valid_counter", "metric_type": "counter"},
			wantCounter: 100,
			wantStatus:  http.StatusOK,
		},
		{
			name:       "wrong method",
			method:     http.MethodPost,
			url:        "/value/counter/valid_counter",
			pathVals:   map[string]string{"metric_name": "valid_counter", "metric_type": "counter"},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "wrong type",
			method:     http.MethodGet,
			url:        "/value/abnormos/qwer",
			pathVals:   map[string]string{"metric_name": "qwer", "metric_type": "abnormos"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found gauge",
			method:     http.MethodGet,
			url:        "/value/gauge/qwer",
			pathVals:   map[string]string{"metric_name": "qwer", "metric_type": "gauge"},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "not found counter",
			method:     http.MethodGet,
			url:        "/value/counter/qwer",
			pathVals:   map[string]string{"metric_name": "qwer", "metric_type": "counter"},
			wantStatus: http.StatusNotFound,
		},
	}

	storage.AddCounter("valid_counter", 100)
	storage.SetGauge("valid_gauge", 100.100)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(tt.method, tt.url, nil)
			for k, v := range tt.pathVals {
				r.SetPathValue(k, v)
			}

			w := httptest.NewRecorder()

			GetMetricHandler(w, r)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
				return
			}

			defer resp.Body.Close()

			if tt.wantStatus != http.StatusOK {
				return
			}

			body, _ := io.ReadAll(resp.Body)
			switch tt.pathVals["metric_type"] {
			case storage.MetricTypeGauge:
				got, _ := strconv.ParseFloat(string(body), 64)
				if got != tt.wantGauge {
					t.Errorf("gauge value = %v, want %v", got, tt.wantGauge)
					return
				}

			case storage.MetricTypeCounter:
				got, _ := strconv.ParseInt(string(body), 10, 64)
				if got != tt.wantCounter {
					t.Errorf("counter value = %v, want %v", got, tt.wantCounter)
					return
				}
			}

		})
	}
}
