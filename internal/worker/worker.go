package worker

import (
	"github.com/jeffreylean/blaster/internal/metrics"
)

type Job interface {
	Do() metrics.Samples
}

type Worker struct {
	ID              int64
	JobChannel      chan Job
	WorkerSharePool chan chan Job
	SampleChannels  chan metrics.Samples
	Rampup          int64
	Exit            chan bool
}

func New(pool chan chan Job, sampCn chan metrics.Samples, rampup int64) *Worker {
	return &Worker{
		JobChannel:      make(chan Job),
		WorkerSharePool: pool,
		SampleChannels:  sampCn,
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
				samp := j.Do()

				w.SampleChannels <- samp
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
