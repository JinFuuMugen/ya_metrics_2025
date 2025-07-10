package storage

import "errors"

const MetricTypeCounter = "counter"
const MetricTypeGauge = "gauge"

type MemStorage struct {
	GaugeMap   map[string]float64
	CounterMap map[string]int64
}

func (ms *MemStorage) GetCounter(k string) (Counter, error) {
	c, exists := ms.CounterMap[k]
	if exists {
		return Counter{Name: k, Type: MetricTypeCounter, Value: c}, nil
	} else {
		return Counter{}, errors.New("missing key")
	}
}

func (ms *MemStorage) GetGauge(k string) (Gauge, error) {
	g, exists := ms.GaugeMap[k]
	if exists {
		return Gauge{Name: k, Type: MetricTypeGauge, Value: g}, nil
	} else {
		return Gauge{}, errors.New("missing key")
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
