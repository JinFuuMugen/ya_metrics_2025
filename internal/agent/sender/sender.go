package sender

import (
	"fmt"
	"net/http"

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

func (v *values) sendMetric(url string) error {
	resp, err := v.client.Post(url, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("cannot send metric: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (v *values) Process(counters []storage.Counter, gauges []storage.Gauge) error {

	for _, c := range counters {

		url := fmt.Sprintf("http://%s/update/%s/%s/%s", v.addr, c.GetType(), c.GetName(), c.GetValueString())

		if err := v.sendMetric(url); err != nil {
			return fmt.Errorf("cannot send counter metric: %w", err)
		}
	}

	for _, g := range gauges {

		url := fmt.Sprintf("http://%s/update/%s/%s/%s", v.addr, g.GetType(), g.GetName(), g.GetValueString())

		if err := v.sendMetric(url); err != nil {
			return fmt.Errorf("cannot send gauge metric: %w", err)
		}
	}

	return nil
}
