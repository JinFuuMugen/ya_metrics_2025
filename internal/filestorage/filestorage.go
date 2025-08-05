package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/config"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/logger"
	models "github.com/JinFuuMugen/ya_metrics_2025/internal/model"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

type Arbitrator struct {
	loader *metricsLoader
	saver  *MetricsSaver
	config *config.ServerConfig
}

func InitArbitrator(cfg *config.ServerConfig) (*Arbitrator, error) {
	loader, err := NewLoader(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create loader: %w", err)
	}

	saver, err := NewSaver(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create saver: %w", err)
	}

	return &Arbitrator{
		loader: loader,
		saver:  saver,
		config: cfg,
	}, nil

}

func (a *Arbitrator) StartArbitrator(ctx context.Context) error {
	if a.config.Restore {
		metrics, err := a.loader.LoadMetrics()
		if err != nil {
			return fmt.Errorf("cannot load metrics from file %s: %w", a.config.FileStoragePath, err)
		}

		for _, m := range metrics {
			switch m.MType {
			case storage.MetricTypeCounter:
				storage.AddCounter(m.ID, *m.Delta)
			case storage.MetricTypeGauge:
				storage.SetGauge(m.ID, *m.Value)
			default:
				return fmt.Errorf("unkown metric type %s :%w", m.MType, err)
			}
		}

		a.loader.Close()
	}

	if a.config.StoreInterval != 0 {
		go a.runFlusher(ctx)
	}

	return nil
}

func (a *Arbitrator) runFlusher(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(a.config.StoreInterval) * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			a.saver.SaveMetrics()

		case <-ctx.Done():
			a.saver.Close()
			return
		}
	}
}

type MetricsSaver struct {
	file *os.File
}

type metricsLoader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewSaver(fname string) (*MetricsSaver, error) {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return &MetricsSaver{file: file}, nil
}

func (ms *MetricsSaver) SaveMetrics() error {

	if err := ms.file.Truncate(0); err != nil {
		return fmt.Errorf("cannot truncate file: %w", err)
	}

	counters := storage.GetCounters()
	gauges := storage.GetGauges()

	metrics := make([]models.Metrics, 0)

	for _, c := range counters {
		metrics = append(metrics, models.Metrics{ID: c.Name, Delta: &c.Value, MType: models.Counter})
	}

	for _, g := range gauges {
		metrics = append(metrics, models.Metrics{ID: g.Name, Value: &g.Value, MType: models.Gauge})
	}

	data, err := json.MarshalIndent(metrics, "", " ")
	if err != nil {
		return fmt.Errorf("cannot marshal metrics: %w", err)
	}

	_, err = ms.file.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write metrics to file: %w", err)
	}

	logger.Infof("data saved: %d metrics", len(metrics))

	return nil
}

func (ms *MetricsSaver) Close() error {
	err := ms.file.Close()
	if err != nil {
		return fmt.Errorf("cannot close file: %w", err)
	}

	return nil
}

func NewLoader(fname string) (*metricsLoader, error) {
	file, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return &metricsLoader{file: file, decoder: json.NewDecoder(file)}, nil
}

func (ml *metricsLoader) LoadMetrics() ([]models.Metrics, error) {
	var metrics []models.Metrics
	if err := ml.decoder.Decode(&metrics); err != nil {
		return nil, fmt.Errorf("cannot decode metrics from file: %w", err)
	}

	return metrics, nil
}

func (ml *metricsLoader) Close() error {
	err := ml.file.Close()
	if err != nil {
		return fmt.Errorf("cannot close file: %w", err)
	}

	return nil
}
