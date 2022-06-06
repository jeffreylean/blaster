package blast

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/jeffreylean/blaster/internal/job"
	"github.com/jeffreylean/blaster/internal/result"
	"github.com/jeffreylean/blaster/internal/scheduler"
	"github.com/kelindar/loader"
)

func Blast(uri, payload string, workers, requests, rampup int64) {
	var err error
	if payload == "" {
		payload, err = getPayload()
		if err != nil {
			log.Printf("blast: unable to read payload, due to %s\n", err.Error())
			return
		}
	}

	start := time.Now()
	// Create scheduler
	wg := new(sync.WaitGroup)
	s := scheduler.New(workers)
	r := new(result.Result)
	r.ResultChannel = s.ResultChannel
	s.Rampup = rampup

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

func getPayload() (string, error) {
	payload := ""
	if e, ok := os.LookupEnv("PAYLOAD"); ok {
		l := loader.New()
		b, err := l.Load(context.Background(), e)
		if err != nil {
			return "", fmt.Errorf("blast: Unable to read payload, due to %s", err.Error())
		}
		payload = string(b)
	}

	// Fill with environment variables
	return payload, nil
}
