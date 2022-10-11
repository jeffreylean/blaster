package job

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	TargetURL string
	Payload   []byte
	WaitGroup *sync.WaitGroup
}

type Response struct {
	WorkerID  int64
	Status    int
	Message   string
	Error     string
	TimeTaken float64
}

func (j *Job) Do() Response {
	defer j.WaitGroup.Done()

	req := new(http.Request)
	client := new(http.Client)
	response := new(Response)

	var err error

	// Default HTTP request set to be POST.
	req, err = http.NewRequest("POST", j.TargetURL, bytes.NewBuffer(j.Payload))
	if err != nil {
		response.Error = err.Error()
		return *response
	}

	// If payload is empty will assume it is GET request.
	if len(j.Payload) == 0 {
		req, err = http.NewRequest("GET", j.TargetURL, nil)
		if err != nil {
			response.Error = err.Error()
			return *response
		}

	}

	// Start hold current time before sending request.
	start := time.Now()
	// Send HTTP request
	resp, err := client.Do(req)
	// Time taken for the request.
	end := time.Since(start)
	if err != nil {
		response.Error = err.Error()
		return *response
	}

	response.Status = resp.StatusCode
	response.Message = resp.Status
	response.TimeTaken = float64(end) / float64(time.Millisecond)
	return *response
}
