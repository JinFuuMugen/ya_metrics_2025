package sender

import "github.com/JinFuuMugen/ya_metrics_2025/internal/storage"

type Sender interface {
	Process([]storage.Counter, []storage.Gauge) error
}
