package main

import (
	"fmt"
	"log"
	"time"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/agent/monitor"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/agent/sender"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/config"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

func main() {

	cfg, err := config.InitAgentConfig()
	if err != nil {
		panic(fmt.Errorf("cannot init agent config: %w", err))
	}

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInerval) * time.Second)

	str := storage.NewStorage()
	snd := sender.NewSender(cfg.Addr)

	m := monitor.NewRuntimeMonitor(str, snd)

	for {
		select {
		case <-pollTicker.C:
			m.CollectRuntimeMetrics()

		case <-reportTicker.C:
			err := m.Dump()
			if err != nil {
				log.Printf("cannot dump metrics: %s\n", err)
			}
		}
	}
}
