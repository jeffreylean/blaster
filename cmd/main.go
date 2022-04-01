package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/jeffreylean/blaster/internal/config"
	"github.com/jeffreylean/blaster/internal/job"
	"github.com/jeffreylean/blaster/internal/scheduler"
)

func main() {
	// Create scheduler
	s := scheduler.New()
	s.Start()

	wg := new(sync.WaitGroup)

	for i := 0; i < int(config.GetConfigInt64("USERS")); i++ {
		wg.Add(1)
		j := new(job.Job)
		j.Payload = `{"schema":"iglu:com.snowplowanalytics.snowplow/payload_data/jsonschema/1-0-4","data":[{"e":"ue","eid":"7aa40ed1-de74-4519-b398-64c276bf1f3c","tv":"js-3.3.0","tna":"ap1","aid":"food","p":"web","cookie":"1","cs":"UTF-8","lang":"en-US","res":"401x746","cd":"30","dtm":"1648649558166","vp":"444x827","ds":"444x827","vid":"5","sid":"f6b6c920-24c9-484d-9e77-8d4b3c86ebf5","duid":"9fec3275-04bc-4ac9-900d-304628d42251","url":"http://localhost:3000/food","ue_px":"eyJzY2hlbWEiOiJpZ2x1OmNvbS5zbm93cGxvd2FuYWx5dGljcy5zbm93cGxvdy91bnN0cnVjdF9ldmVudC9qc29uc2NoZW1hLzEtMC0wIiwiZGF0YSI6eyJzY2hlbWEiOiJpZ2x1OmNvbS5nb29nbGUuYW5hbHl0aWNzLmVuaGFuY2VkLWVjb21tZXJjZS9hY3Rpb24vanNvbnNjaGVtYS8xLTAtMCIsImRhdGEiOnsiYWN0aW9uIjoiY2xpY2sifX19","cx":"eyJzY2hlbWEiOiJpZ2x1OmNvbS5zbm93cGxvd2FuYWx5dGljcy5zbm93cGxvdy9jb250ZXh0cy9qc29uc2NoZW1hLzEtMC0wIiwiZGF0YSI6W3sic2NoZW1hIjoiaWdsdTpjb20uc25vd3Bsb3dhbmFseXRpY3Muc25vd3Bsb3cvd2ViX3BhZ2UvanNvbnNjaGVtYS8xLTAtMCIsImRhdGEiOnsiaWQiOiIwMDI5MTBiYS1jOWViLTQ2NWEtYWI0NC0xYjEzMmQwODA2MDcifX1dfQ","stm":"1648649558168"}]}`
		j.TargetURL = config.GetConfigString("TARGET")
		j.WaitGroup = wg
		s.JobQueue <- j
	}
	wg.Wait()
	fmt.Println("Total Success: ", s.Success)
	fmt.Println("Total Failed: ", s.Fail)
	fmt.Println("Average Time Taken: ", math.Round(s.AverageTimeTaken*100)/100, "ms")
}
