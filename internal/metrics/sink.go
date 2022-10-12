package metrics

var (
	_ Sink = &CounterSink{}
)

type Sink interface {
	Add(s Sample)
	Calculate()
	Get() map[string]float64
}

type CounterSink struct {
	Value float64
}

func (c *CounterSink) Add(s Sample) {
	c.Value += s.Value
}

func (c *CounterSink) Calculate() {
}

// TODO:Data Race happen here, need to fix this
func (c *CounterSink) Get() map[string]float64 {
	return map[string]float64{"count": c.Value}
}
