package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log"

	models "github.com/JinFuuMugen/ya_metrics_2025/internal/model"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

type Sender interface {
	Process([]storage.Counter, []storage.Gauge) error
}

type values struct {
	addr   string
	client *http.Client
}

func NewSender(serverAddr string) *values {
	return &values{serverAddr, &http.Client{}}
}

func (v *values) sendMetric(url string, body io.Reader) error {
	resp, err := v.client.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("cannot send metric: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (v *values) Process(counters []storage.Counter, gauges []storage.Gauge) error {

	for _, c := range counters {

		modeledMetric := models.Metrics{ID: c.Name, MType: storage.MetricTypeCounter, Delta: &c.Value}

		marshMetric, err := json.Marshal(modeledMetric)
		if err != nil {
			log.Printf("cannot marshal metric %s :%s", c.Name, err)
		}

		url := fmt.Sprintf("http://%s/value", v.addr)

		if err := v.sendMetric(url, bytes.NewBuffer(marshMetric)); err != nil {
			return fmt.Errorf("cannot send counter metric: %w", err)
		}
	}

	for _, g := range gauges {

		modeledMetric := models.Metrics{ID: g.Name, MType: storage.MetricTypeGauge, Value: &g.Value}
		marshMetric, err := json.Marshal(modeledMetric)
		if err != nil {
			log.Printf("cannot marshal metric %s :%s", g.Name, err)
		}

		url := fmt.Sprintf("http://%s/value", v.addr)

		if err := v.sendMetric(url, bytes.NewBuffer(marshMetric)); err != nil {
			return fmt.Errorf("cannot send gauge metric: %w", err)
		}
	}

	return nil
}
