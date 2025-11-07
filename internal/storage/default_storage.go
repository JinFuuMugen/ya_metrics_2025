package storage

import (
	"fmt"
)

const MetricTypeCounter = "counter"
const MetricTypeGauge = "gauge"

type MemStorage struct {
	GaugeMap   map[string]float64
	CounterMap map[string]int64
}

func (ms *MemStorage) GetCounter(k string) (Counter, error) {
	c, ok := ms.CounterMap[k]
	if ok {
		return Counter{Name: k, Type: MetricTypeCounter, Value: c}, nil
	} else {
		return Counter{}, fmt.Errorf("missing key: %s", k)
	}
}

func (ms *MemStorage) GetGauge(k string) (Gauge, error) {
	g, ok := ms.GaugeMap[k]
	if ok {
		return Gauge{Name: k, Type: MetricTypeGauge, Value: g}, nil
	} else {
		return Gauge{}, fmt.Errorf("missing key: %s", k)
	}
}

func (ms *MemStorage) GetCounters() []Counter {

	var counters []Counter

	for k, v := range ms.CounterMap {
		counters = append(counters, Counter{Name: k, Type: MetricTypeCounter, Value: v})
	}
	return counters
}

func (ms *MemStorage) GetGauges() []Gauge {

	var gauges []Gauge

	for k, v := range ms.GaugeMap {
		gauges = append(gauges, Gauge{Name: k, Type: MetricTypeGauge, Value: v})
	}
	return gauges
}

func (ms *MemStorage) SetGauge(k string, v float64) {
	ms.GaugeMap[k] = v
}

func (ms *MemStorage) AddCounter(k string, v int64) {
	ms.CounterMap[k] += v
}

func (ms *MemStorage) ResetCounters() {
	for k := range ms.CounterMap {
		ms.CounterMap[k] = 0
	}
}
