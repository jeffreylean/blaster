package scheduler

import (
	"fmt"
	"time"

	"github.com/jeffreylean/blaster/internal/metrics"
	"github.com/jeffreylean/blaster/internal/worker"
)

type Scheduler struct {
	WorkerPool    chan chan worker.Job
	JobQueue      chan worker.Job
	MaxWorker     int64
	Rampup        int64 // In seconds
	StopWorker    bool
	SampleChannel chan metrics.Samples
}

func New(workers, rampup int64) *Scheduler {
	return &Scheduler{
		WorkerPool:    make(chan chan worker.Job, workers),
		MaxWorker:     workers,
		JobQueue:      make(chan worker.Job),
		SampleChannel: make(chan metrics.Samples),
		StopWorker:    false,
		Rampup:        rampup,
	}
}

func (s *Scheduler) Start() {
	go func() {
		fmt.Println("Starting all workers......")
		s.Dispatch()

		for i := int64(0); i < s.MaxWorker; i++ {
			if s.StopWorker {
				break
			}
			go func() {
				w := worker.New(s.WorkerPool, s.SampleChannel, s.Rampup)
				w.Work()
			}()
			var waitDuration float32 = float32(s.Rampup) / float32(s.MaxWorker)
			time.Sleep(time.Duration(int(1000*waitDuration)) * time.Millisecond)
		}
	}()
}

func (s *Scheduler) Dispatch() {
	go func() {
		for {
			select {
			case job := <-s.JobQueue:
				jobChannel := <-s.WorkerPool
				jobChannel <- job
			}
		}
	}()
}
