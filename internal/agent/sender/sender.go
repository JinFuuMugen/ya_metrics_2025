package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"log"

	models "github.com/JinFuuMugen/ya_metrics_2025/internal/model"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
	"github.com/go-resty/resty/v2"
)

type Sender interface {
	Process([]storage.Counter, []storage.Gauge) error
}

type values struct {
	addr   string
	client *resty.Client
}

func NewSender(serverAddr string) *values {
	return &values{serverAddr, resty.New()}
}

func (v *values) sendMetric(url string, body io.Reader) error {

	_, err := v.client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(body).Post(url)
	if err != nil {
		return fmt.Errorf("cannot send metric: %w", err)
	}

	return nil
}

func (v *values) Compress(data []byte) ([]byte, error) {

	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("cannot write data to compress: %w", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("cannot compress data: %w", err)
	}

	return b.Bytes(), nil
}

func (v *values) Process(counters []storage.Counter, gauges []storage.Gauge) error {

	for _, c := range counters {

		modeledMetric := models.Metrics{ID: c.Name, MType: storage.MetricTypeCounter, Delta: &c.Value}

		marshMetric, err := json.Marshal(modeledMetric)
		if err != nil {
			log.Printf("cannot marshal metric %s :%s", c.Name, err)
		}

		url := fmt.Sprintf("http://%s/update", v.addr)

		compressedData, err := v.Compress(marshMetric)
		if err != nil {
			return fmt.Errorf("cannot compress data: %w", err)
		}

		if err := v.sendMetric(url, bytes.NewBuffer(compressedData)); err != nil {
			return fmt.Errorf("cannot send counter metric: %w", err)
		}
	}

	for _, g := range gauges {

		modeledMetric := models.Metrics{ID: g.Name, MType: storage.MetricTypeGauge, Value: &g.Value}
		marshMetric, err := json.Marshal(modeledMetric)
		if err != nil {
			log.Printf("cannot marshal metric %s :%s", g.Name, err)
		}

		url := fmt.Sprintf("http://%s/update", v.addr)

		compressedData, err := v.Compress(marshMetric)
		if err != nil {
			return fmt.Errorf("cannot compress data: %w", err)
		}

		if err := v.sendMetric(url, bytes.NewBuffer(compressedData)); err != nil {
			return fmt.Errorf("cannot send gauge metric: %w", err)
		}
	}

	return nil
}
