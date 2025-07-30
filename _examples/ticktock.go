package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/ticktock"
	"github.com/hootuu/hyle/data/idx"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		worker := ticktock.NewWorker("main_tt", zplt.HelixRdsCache())
		worker.RegisterJobHandlerFunc("TIMEOUT", func(ctx context.Context, job *ticktock.Job) error {
			fmt.Println(">>>>>>[", time.Now().Format("2006-01-02 15:04:05"), "]TIMEOUT", job.Type, string(job.Payload))
			return nil
		})
		postman := ticktock.NewPostman("postman_tt", zplt.HelixRdsCache())

		for i := 0; i < 0; i++ {
			err := postman.Send(context.Background(), &ticktock.NowJob{
				Type:      "TIMEOUT",
				ID:        fmt.Sprintf("timeout_%d_%d", i, time.Now().UnixMilli()),
				Payload:   []byte(fmt.Sprintf("timeout_%d_%d", i, time.Now().UnixMilli())),
				UniqueTTL: 0,
			})
			if err != nil {
				panic(err)
			}
		}

		for i := 0; i < 0; i++ {
			du := time.Duration(i+1) * 1 * time.Second
			willRun := time.Now().Add(du)
			willRunTimeStr := willRun.Format("2006-01-02 15:04:05")
			err := postman.Send(context.Background(), &ticktock.DelayJob{
				Type:      "TIMEOUT",
				ID:        fmt.Sprintf("delay_timeout_%d_%s", i, willRunTimeStr),
				Payload:   []byte(fmt.Sprintf("delay_timeout_%d_%s", i, willRunTimeStr)),
				UniqueTTL: 0,
				Delay:     du,
			})
			if err != nil {
				panic(err)
			}
		}

		for i := 0; i < 0; i++ {
			//du := time.Duration(i+1) * 1 * time.Second
			willRun := time.Now().Truncate(time.Minute).Add(time.Minute)
			willRunTimeStr := willRun.Format("2006-01-02 15:04:05")
			err := postman.Send(context.Background(), &ticktock.RunAtJob{
				Type:      "TIMEOUT",
				ID:        fmt.Sprintf("run_at_timeout_%d_%s", i, willRunTimeStr),
				Payload:   []byte(fmt.Sprintf("run_at_timeout_%d_%s", i, willRunTimeStr)),
				UniqueTTL: 0,
				RunAt:     willRun,
			})
			if err != nil {
				panic(err)
			}
		}

		go func() {
			for i := 0; i < 1; i++ {
				du := time.Duration(i+1) * 1 * time.Second
				nxtTime, _ := ticktock.Expression("* * * * *").Next(time.Now())
				fmt.Println("nxtTimenxtTimenxtTime", nxtTime)
				willRun := time.Now().Add(du)
				willRunTimeStr := willRun.Format("2006-01-02 15:04:05")
				err := postman.Send(context.Background(), &ticktock.PeriodicJob{
					Type:         "TIMEOUT",
					Payload:      []byte(fmt.Sprintf("period_timeout_%d_%s", i, willRunTimeStr)),
					UniqueTTL:    0,
					JobTplID:     "TPL_" + idx.New(),
					TplUniqueTTL: 0,
					Sequence:     0,
					Expression:   "* * * * *",
					Current:      nxtTime,
				})
				if err != nil {
					panic(err)
				}
				//time.Sleep(1 * time.Second)
			}
		}()
		fmt.Println("add job success")

		time.Sleep(1 * time.Hour)
	})
	helix.Startup()
}
