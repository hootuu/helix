package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
)

func main() {
	var uidGenerator hnid.Generator
	func() {
		h := helix.BuildHelix("biztest", func() (context.Context, error) {
			var err error
			uidGenerator, err = hnid.NewGenerator("examples.uid",
				hnid.NewOptions(2, 99).
					SetTimestamp(hnid.Second, true).
					SetAutoInc(8, 1, 99999999, 10000))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return nil, nil
		}, func(ctx context.Context) {

		})
		helix.Use(h)
	}()
	helix.AfterStartup(func() {
		for i := 0; i < 20; i++ {
			id := uidGenerator.Next()
			fmt.Println(hjson.MustToString(id))
			fmt.Println(id.ToString())
		}
	})
	helix.Startup()
}
