package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/sync/croupier"
	"github.com/hootuu/hyle/hypes/collar"
)

func main() {
	helix.AfterStartup(func() {
		cp := croupier.Light(collar.Build("TEST", "11199"))
		err := cp.Publish(context.Background(),
			13,
			0,
		)
		if err != nil {
			panic(err)
		}
		for i := 0; i < 18; i++ {
			b, err := cp.Allow(context.Background(), func() error {
				if i == 11 {
					return errors.New("test error")
				} else {
					//fmt.Println("allow ok do")
				}
				return nil
			})
			fmt.Println("[", i+1, "]", b, err)
		}
	})
	helix.Startup()
}
