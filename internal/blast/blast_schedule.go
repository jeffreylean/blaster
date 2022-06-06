package blast

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffreylean/blaster/internal/config"
)

func ScheduleBlast(targetUrl, payload string, workers, rampup uint32) error {
	c := config.LoadOrPanic()

	nextTime := time.Now()
	end := time.Now().AddDate(0, 0, 1)
	for nextTime.Format("2006-01-02") != end.Format("2006-01-02") {
		fmt.Println(time.Now(), " blasting....")
		req, err := readSchedule(*c)
		if err != nil {
			log.Printf("schedule blast: unable to schedule, due to %s\n", err.Error())
		}
		go Blast(targetUrl, payload, int64(workers), int64(req), int64(rampup))
		nextTime = nextTime.Add(time.Second * time.Duration(c.Iteration))
		time.Sleep(time.Second * time.Duration(c.Iteration))
	}
	return nil
}

func readSchedule(c config.Config) (int, error) {
	currHour := time.Now().Hour()
	requests := 0
	for _, each := range c.Schedule {
		if each.Hour == currHour {
			requests = int((float64(c.Iteration) / 3600) * float64(each.Requests))
			break
		}
	}
	return requests, nil
}
