package job

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/jeffreylean/blaster/internal/metrics"
)

type Job struct {
	TargetURL string
	Payload   []byte
	WaitGroup *sync.WaitGroup
	Metrics   metrics.Metrics
}

type Response struct {
	Status    int
	Message   string
	Error     string
	TimeTaken float64
}

// Convert response data into metrics sample for metrics aggregation
func (r Response) toSamples(m metrics.Metrics) metrics.Samples {
	s := make(metrics.Samples, 0)
	if r.Status >= 200 && r.Status <= 299 && r.Error == "" {
		success := m[metrics.SUCCESS]
		samp := metrics.Sample{
			Metric: &success,
			Value:  1,
		}
		s = append(s, samp)
	}

	fail := m[metrics.FAIL]
	if r.Status > 400 && r.Error != "" {
		samp := metrics.Sample{
			Metric: &fail,
			Value:  1,
		}
		s = append(s, samp)
	}

	return s
}

func (j *Job) Do() metrics.Samples {
	defer j.WaitGroup.Done()

	req := new(http.Request)
	client := new(http.Client)
	response := new(Response)

	var err error

	// Default HTTP request set to be POST
	req, err = http.NewRequest("POST", j.TargetURL, bytes.NewBuffer(j.Payload))
	if err != nil {
		response.Error = err.Error()
		return response.toSamples(j.Metrics)
	}

	// If payload is empty will assume it is GET request
	if len(j.Payload) == 0 {
		req, err = http.NewRequest("GET", j.TargetURL, nil)
		if err != nil {
			response.Error = err.Error()
			return response.toSamples(j.Metrics)
		}

	}

	// Start hold current time before sending request
	start := time.Now()
	// Send HTTP request
	resp, err := client.Do(req)
	// Time taken for the request
	end := time.Since(start)
	if err != nil {
		response.Error = err.Error()
		return response.toSamples(j.Metrics)
	}

	response.Status = resp.StatusCode
	response.Message = resp.Status
	response.TimeTaken = float64(end) / float64(time.Millisecond)

	// Convert into Sample type
	return response.toSamples(j.Metrics)
}
