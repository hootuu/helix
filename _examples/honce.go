package main

import (
	"fmt"
	"github.com/hootuu/helix/components/honce"
	"github.com/hootuu/helix/helix"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		result := ""
		err := honce.Do(fmt.Sprintf("examples.do.%d", time.Now().Unix()), func() error {
			result = "one"
			return nil
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result) //should be one
		err = honce.Do(fmt.Sprintf("examples.do.%d", time.Now().Unix()), func() error {
			result = "two"
			return nil
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result) // should be one
	})
	helix.Startup()
}
