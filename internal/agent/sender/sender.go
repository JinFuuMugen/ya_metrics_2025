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

func NewSender() *values {
	return &values{"http://localhost:8080", &http.Client{}} //make dynamical
}

func (v *values) Process(counters []storage.Counter, gauges []storage.Gauge) error {

	for _, c := range counters {

		url := fmt.Sprintf("%s/update/%s/%s/%s", v.addr, c.GetType(), c.GetName(), c.GetValueString())

		resp, err := v.client.Post(url, "text/plain", nil)
		if err != nil {
			return fmt.Errorf("cannot send counter metric: %w", err)
		}

		defer resp.Body.Close()

	}

	for _, g := range gauges {

		url := fmt.Sprintf("%s/update/%s/%s/%s", v.addr, g.GetType(), g.GetName(), g.GetValueString())

		resp, err := v.client.Post(url, "text/plain", nil)
		if err != nil {
			return fmt.Errorf("cannot send gauge metric: %w", err)
		}

		defer resp.Body.Close()
	}

	return nil
}
