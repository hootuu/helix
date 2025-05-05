package main

import (
	"fmt"
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
)

func main() {
	helix.OnStartupSuccess(func() {
		uidGenerator, err := hnid.NewGenerator("examples.uid",
			hnid.NewOptions(2, 99).
				SetTimestamp(hnid.Second, true).
				SetAutoInc(8, 1, 99999999, 10000))
		if err != nil {
			fmt.Println(err)
			return
		}
		for i := 0; i < 20; i++ {
			id := uidGenerator.Next()
			fmt.Println(hjson.MustToString(id))
			fmt.Println(id.ToString())
		}
	})
	helix.Startup()
}
