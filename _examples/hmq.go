package main

import (
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/helix/unicom/hmq/hnsq"
	"sync"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		nsqMQ := hnsq.NewNsqMQ()
		mq := hmq.NewMQ("main_mq", nsqMQ)
		thisP := hmq.NewProducer("a", "topic_a").With(nsqMQ.NewProducer())
		err := mq.RegisterProducer(thisP)
		if err != nil {
			panic(err)
		}
		thisC := hmq.NewConsumer("b", "topic_a", "1").
			With(nsqMQ.NewConsumer()).WithHandler(func(msg *hmq.Message) error {
			fmt.Println("====>>>>>>>", string(msg.Payload[:]))
			return nil
		})
		err = mq.RegisterConsumer(thisC)
		if err != nil {
			panic(err)
		}
		idx := uint64(0)
		mu := sync.Mutex{}
		for i := 0; i < 1000; i++ {
			go func() {
				for j := 0; j < 1000; j++ {
					var cur uint64
					mu.Lock()
					cur = idx
					idx++
					mu.Unlock()
					msg := fmt.Sprintf("hello world %d", cur)
					err = thisP.Publish([]byte(msg))
					if err != nil {
						fmt.Println(err)
					}
				}
			}()
		}

		time.Sleep(9 * time.Hour)
	})
	helix.Startup()
}
