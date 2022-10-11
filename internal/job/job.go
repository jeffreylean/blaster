package job

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	TargetURL string
	Payload   string
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
	start := time.Now()
	client := new(http.Client)
	response := new(Response)

	b := []byte(j.Payload)

	req, err := http.NewRequest("POST", j.TargetURL, bytes.NewBuffer(b))
	end := time.Since(start)
	if err != nil {
		response.Error = err.Error()
		return *response
	}

	resp, err := client.Do(req)
	if err != nil {
		response.Error = err.Error()
		return *response
	}

	response.Status = resp.StatusCode
	response.Message = resp.Status
	response.TimeTaken = float64(end) / float64(time.Millisecond)
	return *response
}
