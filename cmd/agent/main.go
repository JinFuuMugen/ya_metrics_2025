package main

import (
	"fmt"
	"time"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

func main() {

	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	storage.InitStorage()

	for {
		select {
		case <-pollTicker.C:
			//do polling
			fmt.Println("poll")

		case <-reportTicker.C:
			//do reporting
			fmt.Println("report")
		}
	}
}
