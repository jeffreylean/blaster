package blast

import (
	"fmt"
	"sync"
	"time"

	"github.com/jeffreylean/blaster/internal/job"
	"github.com/jeffreylean/blaster/internal/metrics"
	"github.com/jeffreylean/blaster/internal/scheduler"
)

func Blast(uri string, payload []byte, workers, requests, rampup int64) {
	start := time.Now()
	wg := new(sync.WaitGroup)

	// Create scheduler
	s := scheduler.New(workers, rampup)

	// Create metrics
	m := metrics.New(s.SampleChannel, requests)
	m.AggregateMetrics()

	s.Start()
	for i := 0; i < int(requests); i++ {
		wg.Add(1)
		j := new(job.Job)
		j.Payload = payload
		j.TargetURL = uri
		j.WaitGroup = wg
		j.Metrics = m.Metrics
		s.JobQueue <- j
	}

	// Wait for all request finish handle by workers
	wg.Wait()
	// Wait for metrics aggregation
	m.Wg.Wait()

	s.StopWorker = true
	end := time.Since(start)
	m.Summary()
	//fmt.Println("Total Success: ", r.Success)
	//fmt.Println("Total Failed: ", r.Fail)
	//fmt.Println("Average Time Taken: ", math.Round(r.AverageTimeTaken*1000)/1000, "ms")
	fmt.Println("Total Time Taken: ", end.Seconds(), "s")
}
