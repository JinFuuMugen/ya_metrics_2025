package handler

import (
	"net/http"
	"net/http/httptest"
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
			}

			if tt.wantGauge != 0 {
				got, err := storage.GetGauge("GaugeMetr")
				if got.GetValue() == 0 || got.GetValue() != tt.wantGauge || err != nil {
					t.Errorf("gauge = %v, want %v", got.GetValue(), tt.wantGauge)
				}
			}

			if tt.wantCounter != 0 {
				got, err := storage.GetCounter("CounterMetr")
				if got == got.GetValue() || got.GetValue() != tt.wantCounter || err != nil {
					t.Errorf("counter = %v, want %v", got, tt.wantGauge)
				}
			}
		})
	}
}
