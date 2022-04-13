package scheduler

import (
	"fmt"
	"time"

	"github.com/jeffreylean/blaster/internal/worker"
)

type Scheduler struct {
	WorkerPool    chan chan worker.Job
	JobQueue      chan worker.Job
	MaxWorker     int64
	ResultChannel chan any
	Rampup        int64 // In millisecond
}

func New(workers int64) *Scheduler {
	max := workers
	return &Scheduler{
		WorkerPool:    make(chan chan worker.Job, max),
		MaxWorker:     max,
		JobQueue:      make(chan worker.Job),
		ResultChannel: make(chan any),
	}
}

func (s *Scheduler) Start() {
	go func() {
		fmt.Println("Starting all workers......")
		s.Dispatch()
		for i := int64(0); i < s.MaxWorker; i++ {
			go func(id int64) {
				fmt.Println("Worker ", id, " created...")
				w := worker.New(s.WorkerPool, s.ResultChannel, id, s.Rampup)
				w.Work()
			}(i)
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
