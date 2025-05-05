package main

import (
	"fmt"
	"github.com/hootuu/helix/components/hidem"
	"github.com/hootuu/helix/helix"
	"time"
)

func main() {
	helix.OnStartupSuccess(func() {
		idemFactory, err := hidem.NewDbFactory("example_order", 30*time.Minute, 5*time.Minute)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := idemFactory.Check("abc")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(b)
	})
	helix.Startup()
}
