package main

import (
	"fmt"
	"github.com/hootuu/helix/components/htick"
	"github.com/hootuu/helix/helix"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		//err := htick.Schedule(&htick.Job{
		//	Expression: "*/32 * * * * ?",
		//	Topic:      "HELLO_WORLD",
		//	Payload:    []byte("HELLO_WORLD"),
		//})
		//fmt.Println(err)
		err := htick.Once(time.Now().Add(5*time.Second), "hello_World", []byte("a"))
		fmt.Println(err)
	})
	helix.Startup()
}
