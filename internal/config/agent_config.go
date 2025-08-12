package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	Addr          string `env:"ADDRESS"`
	PollInterval  int    `env:"POLL_INTERVAL"`
	ReportInerval int    `env:"REPORT_INTERVAL"`
}

func InitAgentConfig() (*AgentConfig, error) {

	agentConfig := new(AgentConfig)

	flagAddr := flag.String("a", "localhost:8080", "Metrics server address")
	flagPollInterval := flag.Int("p", 2, "Metrics polling interval")
	flagReportInterval := flag.Int("r", 10, "Metrics report interval")

	flag.Parse()

	agentConfig.Addr = *flagAddr
	agentConfig.PollInterval = *flagPollInterval
	agentConfig.ReportInerval = *flagReportInterval

	err := env.Parse(agentConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot parse env variables: %w", err)
	}

	return agentConfig, nil
}
