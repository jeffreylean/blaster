package result

import (
	"fmt"

	"github.com/jeffreylean/blaster/internal/job"
)

type Result struct {
	Success          int64
	Fail             int64
	AverageTimeTaken float64
	ResultChannel    chan any
}

func (r *Result) PrintResult() {
	go func() {
		avg := 0.0
		total := 0.0
		for response := range r.ResultChannel {
			resp := response.(job.Response)
			if resp.Error == "" && resp.Status >= 200 && resp.Status <= 299 {
				r.Success += 1
				total += resp.TimeTaken
				avg = total / float64(r.Success)
				r.AverageTimeTaken = avg
				fmt.Println(fmt.Sprintf("Success: worker %d used %.3f ms to complete the job...", resp.WorkerID, resp.TimeTaken))
			} else {
				r.Fail += 1
				fmt.Println(fmt.Sprintf("Failed: worker %d failed the job due to %s", resp.WorkerID, resp.Error))
			}
		}

	}()
}
