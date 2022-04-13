package worker

import (
	"github.com/jeffreylean/blaster/internal/job"
)

type Job interface {
	Do() job.Response
}

type Worker struct {
	ID              int64
	JobChannel      chan Job
	WorkerSharePool chan chan Job
	ResultChannel   chan any
	Rampup          int64
	Exit            chan bool
}

func New(pool chan chan Job, result chan any, id int64, rampup int64) *Worker {
	return &Worker{
		ID:              id,
		JobChannel:      make(chan Job),
		WorkerSharePool: pool,
		ResultChannel:   result,
		Exit:            make(chan bool),
		Rampup:          rampup,
	}
}

func (w *Worker) Work() {
	go func() {
		for {
			// Register worker to global worker share pool
			w.WorkerSharePool <- w.JobChannel
			// Wait for job
			select {
			case j := <-w.JobChannel:
				resp := j.Do()
				resp.WorkerID = w.ID
				w.ResultChannel <- resp
			case <-w.Exit:
				return
			}
		}
	}()
}

func (w *Worker) Rest() {
	go func() {
		w.Exit <- true
	}()
}
