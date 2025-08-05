package filestorage

import (
	"encoding/json"
	"fmt"
	"os"

	models "github.com/JinFuuMugen/ya_metrics_2025/internal/model"
)

// func LoadMetrics(filepath string) error {
// 	metricsData, err :=
// }

type MetricsSaver struct {
	file *os.File
}

type MetricsLoader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewSaver(fname string) (*MetricsSaver, error) {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return &MetricsSaver{file: file}, nil
}

func (ms *MetricsSaver) SaveMetric(metric models.Metrics) error {
	data, err := json.MarshalIndent(metric, "", " ")
	if err != nil {
		return fmt.Errorf("cannot marshal metrics: %w", err)
	}

	_, err = ms.file.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write metrics to file: %w", err)
	}

	return nil
}

func (ms *MetricsSaver) Close() error {
	err := ms.file.Close()
	if err != nil {
		return fmt.Errorf("cannot close file: %w", err)
	}

	return nil
}

func NewLoader(fname string) (*MetricsLoader, error) {
	file, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return &MetricsLoader{file: file, decoder: json.NewDecoder(file)}, nil
}

func (ml *MetricsLoader) LoadMetrics() ([]models.Metrics, error) {
	var metrics []models.Metrics
	if err := ml.decoder.Decode(metrics); err != nil {
		return nil, fmt.Errorf("cannot decode metrics from file: %w", err)
	}

	return metrics, nil
}
