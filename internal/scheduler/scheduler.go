package scheduler

import (
	"fmt"

	"github.com/jeffreylean/blaster/internal/config"
	"github.com/jeffreylean/blaster/internal/job"
	"github.com/jeffreylean/blaster/internal/worker"
)

type Scheduler struct {
	WorkerPool       chan chan worker.Job
	JobQueue         chan worker.Job
	MaxWorker        int64
	ResultChannel    chan any
	Success          int64
	Fail             int64
	AverageTimeTaken float64
}

func New() *Scheduler {
	max := config.GetConfigInt64("WORKERS")
	return &Scheduler{
		WorkerPool:    make(chan chan worker.Job, max),
		MaxWorker:     max,
		JobQueue:      make(chan worker.Job),
		ResultChannel: make(chan any),
	}
}

func (s *Scheduler) Start() {
	fmt.Println("Starting all workers......")
	for i := int64(0); i < s.MaxWorker; i++ {
		go func(id int64) {
			w := worker.New(s.WorkerPool, s.ResultChannel, id)
			w.Work()
		}(i)
	}

	s.Dispatch()
}

func (s *Scheduler) Dispatch() {
	go func() {
		avg := 0.0
		total := 0.0
		for {
			select {
			case job := <-s.JobQueue:
				jobChannel := <-s.WorkerPool
				jobChannel <- job
			case response := <-s.ResultChannel:
				resp := response.(job.Response)
				if resp.Error == "" && resp.Status >= 200 && resp.Status <= 299 {
					s.Success += 1
					total += resp.TimeTaken
					avg = total / float64(s.Success)
					s.AverageTimeTaken = avg
					fmt.Println(fmt.Sprintf("Success: worker %d used %.2f ms to complete the job...", resp.WorkerID, resp.TimeTaken))
				} else {
					s.Fail += 1
					fmt.Println(fmt.Sprintf("Failed: worker %d failed the job due to %s", resp.WorkerID, resp.Error))
				}
			}
		}
	}()
}
