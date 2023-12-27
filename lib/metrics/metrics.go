package metrics

import (
	"fmt"

	lib "github.com/avito-tech/avito-pixel/lib"
)

type Metrics struct {
	baseMetrics lib.Metrics
}

func NewMetrics(metrics lib.Metrics) Metrics {
	return Metrics{baseMetrics: metrics}
}

func getMetricsKey(key string) string {
	return fmt.Sprintf("avito_pixel.%s", key)
}

func (m *Metrics) Increment(key string) {
	m.baseMetrics.Increment(getMetricsKey(key))
}

func (m *Metrics) Timer(key string) func() {
	return m.baseMetrics.Timer(key)
}
