package blast

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/jeffreylean/blaster/internal/job"
	"github.com/jeffreylean/blaster/internal/result"
	"github.com/jeffreylean/blaster/internal/scheduler"
)

func Blast(uri string, payload []byte, workers, requests, rampup int64) {
	start := time.Now()
	wg := new(sync.WaitGroup)

	// Create scheduler
	s := scheduler.New(workers)
	s.Rampup = rampup

	r := new(result.Result)
	r.ResultChannel = s.ResultChannel

	s.Start()
	r.PrintResult()

	for i := 0; i < int(requests); i++ {
		wg.Add(1)
		j := new(job.Job)
		j.Payload = payload
		j.TargetURL = uri
		j.WaitGroup = wg
		s.JobQueue <- j
	}
	wg.Wait()

	s.StopWorker = true
	end := time.Since(start)
	// Add some buffer time for result to finish aggregate
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("Total Success: ", r.Success)
	fmt.Println("Total Failed: ", r.Fail)
	fmt.Println("Average Time Taken: ", math.Round(r.AverageTimeTaken*1000)/1000, "ms")
	fmt.Println("Total Time Taken: ", end.Seconds(), "s")
}
