package storage

const MetricTypeCounter = "counter"
const MetricTypeGauge = "gauge"

type (
	Gauge struct {
		Name  string
		Type  string
		Value float64
	}

	Counter struct {
		Name  string
		Type  string
		Value int64
	}

	MemStorage struct {
		GaugeMap   map[string]float64
		CounterMap map[string]int64
	}
)

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
