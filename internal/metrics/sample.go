package metrics

// Unit of metric measurement
type Sample struct {
	Metric *Metric
	Value  float64
}

// Samples consists of many sample
type Samples []Sample
