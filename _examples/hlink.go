package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/hlink"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {

		for i := 0; i < 1000; i++ {

			ctx := context.Background()
			inviteCode, err := hlink.Generate(
				ctx,
				"NINEORA_USER_INVITE",
				collar.Build("USER", "uid_0000008_"+cast.ToString(time.Now().UnixNano())),
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("invite code:", inviteCode)

			err = hlink.Bind(
				inviteCode,
				"FRIEND",
				collar.Build("USER", "uid_0000002"),
			)
			if err != nil {
				fmt.Println(err)
				return
			}

		}

	})

	helix.Startup()
}
