package storage

var defaultStorage MemStorage

func InitStorage() {
	defaultStorage.CounterMap = make(map[string]int64)
	defaultStorage.GaugeMap = make(map[string]float64)
}

func SetGauge(k string, v float64) {
	defaultStorage.SetGauge(k, v)
}

func AddCounter(k string, v int64) {
	defaultStorage.AddCounter(k, v)
}

func GetCounters() []Counter {
	return defaultStorage.GetCounters()
}

func GetGauges() []Gauge {
	return defaultStorage.GetGauges()
}
