package metrics

import (
	"fmt"
	"sync"
)

const (
	CounterType = MetricType(iota)
	GaugeType
	RateType
	TrendType
)

const (
	SUCCESS = MetricName("success")
	FAIL    = MetricName("fail")
)

// Type of Metric
type MetricType int

// Name of Metric
type MetricName string

// Map of metric with name as key
type Metrics map[MetricName]Metric

// Registry which run metrics
type MetricRegistry struct {
	// All defined metrics
	Metrics map[MetricName]Metric
	// Chan where response is sent, source of aggregation
	SampleChan chan Samples
	// Waitgroup which indicate completion aggregation
	Wg *sync.WaitGroup
	// Number of requests (use this value to know exactly how many samples we should expect from sample channel)
	requests int64
}

// Print out summary of metrics
func (m MetricRegistry) Summary() {
	for _, each := range m.Metrics {
		data := each.Sink.Get()

		for k, v := range data {
			fmt.Println(k, ": ", v)
		}
	}
}

// Aggregate metrics data from sample received from sample channel
func (m *MetricRegistry) AggregateMetrics() {
	go func() {
		defer m.Wg.Done()
		for x := 0; int64(x) < m.requests; x++ {
			samp := <-m.SampleChan
			for _, each := range samp {
				each.Metric.Sink.Add(each)
			}
		}
	}()
}

type Metric struct {
	// Name of the metric
	Name MetricName
	// Type of the metric
	MetricType MetricType
	// Sink of the metric
	Sink Sink
}

func New(sampCn chan Samples, requests int64) MetricRegistry {
	m := make(map[MetricName]Metric, 0)

	// Success count metrics
	success := Metric{
		Name:       SUCCESS,
		MetricType: CounterType,
		Sink:       &CounterSink{},
	}
	m[success.Name] = success

	// Fail count metrics
	fail := Metric{
		Name:       FAIL,
		MetricType: CounterType,
		Sink:       &CounterSink{},
	}
	m[success.Name] = fail

	mr := MetricRegistry{
		Metrics:    m,
		SampleChan: sampCn,
		Wg:         new(sync.WaitGroup),
		requests:   requests,
	}

	// Add count
	mr.Wg.Add(1)
	return mr
}
