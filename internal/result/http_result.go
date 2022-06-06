package result

import (
	"errors"
	"fmt"
	"os"
	"time"

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
				//				fmt.Println(fmt.Sprintf("Success: worker %d used %.3f ms to complete the job...", resp.WorkerID, resp.TimeTaken))
			} else {
				r.Fail += 1
				// Check if error file exists
				_, err := os.Stat(fmt.Sprintf("error_%s.txt", time.Now().Format("2006-01-02")))
				if err != nil {
					// Create new file
					if errors.Is(err, os.ErrNotExist) {
						os.WriteFile(fmt.Sprintf("error_%s.txt", time.Now().Format("2006-01-02")), []byte(fmt.Sprintf("%s Error: %s\n", time.Now().String(), resp.Error)), 0644)
					}
				} else {
					// Overwriting the existing file
					file, _ := os.OpenFile(fmt.Sprintf("error_%s.txt", time.Now().Format("2006-01-02")), os.O_APPEND|os.O_WRONLY, 0644)
					file.Write([]byte(fmt.Sprintf("%s Error: %s\n", time.Now().String(), resp.Error)))
					file.Close()
				}
				//				fmt.Println(fmt.Sprintf("Failed: worker %d failed the job due to %s", resp.WorkerID, resp.Error))
			}
		}
	}()
}
