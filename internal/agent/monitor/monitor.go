package monitor

import (
	"github.com/JinFuuMugen/ya_metrics_2025/internal/agent/sender"
)

type Monitor interface {
	Collect()
	Dump() error
	SetProcessor(p sender.Sender)
}

type RuntimeMonitor interface {
	Monitor
	CollectRuntimeMetrics()
}

type GopsutilMonitor interface {
	Monitor
	CollectGopsutil()
}
