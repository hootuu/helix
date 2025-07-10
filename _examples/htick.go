package main

import (
	"fmt"
	"github.com/hootuu/helix/components/htick"
	"github.com/hootuu/helix/helix"
)

func main() {
	helix.AfterStartup(func() {
		err := htick.Schedule(&htick.Job{
			Expression: "*/32 * * * * ?",
			Topic:      "HELLO_WORLD",
			Payload:    []byte("HELLO_WORLD"),
		})
		fmt.Println(err)
	})
	helix.Startup()
}
