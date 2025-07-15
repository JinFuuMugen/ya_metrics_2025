package config

import "flag"

type AgentConfig struct {
	Addr          string
	PollInterval  int
	ReportInerval int
}

func InitAgentConfig() *AgentConfig {

	addr := flag.String("a", "localhost:8080", "Metrics server address")
	pollInterval := flag.Int("p", 2, "Metrics polling interval")
	reportInterval := flag.Int("r", 10, "Metrics report interval")

	flag.Parse()

	return &AgentConfig{Addr: *addr, ReportInerval: *reportInterval, PollInterval: *pollInterval}
}
