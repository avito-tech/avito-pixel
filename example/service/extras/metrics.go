package extras

type Metrics struct{}

func NewMetrics() Metrics {
	return Metrics{}
}

func (m Metrics) Increment(key string) {}
func (m Metrics) Timer(key string) func() {
	return func() {}
}
