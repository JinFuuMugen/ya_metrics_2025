package main

import (
	"time"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/agent/monitor"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/agent/sender"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

func main() {

	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	str := storage.NewStorage()
	snd := sender.NewSender()

	m := monitor.NewRuntimeMonitor(str, snd)

	for {
		select {
		case <-pollTicker.C:
			m.CollectRuntimeMetrics()

		case <-reportTicker.C:
			err := m.Dump()
			if err != nil {
				panic(err)
			}
		}
	}
}
