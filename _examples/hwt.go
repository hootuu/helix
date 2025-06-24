package main

import (
	"fmt"
	"github.com/hootuu/helix/components/hwt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
)

func main() {

	helix.AfterStartup(func() {
		hwt, err := hwt.NewHwt("example_hwt")
		if err != nil {
			fmt.Println(err)
			return
		}
		jwtToken, err := hwt.RefreshIssuing("example_identity")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(hjson.MustToString(jwtToken))
		tJwtToken, err := hwt.TokenIssuing("example_identity", jwtToken.Refresh)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(hjson.MustToString(tJwtToken))
	})
	helix.Startup()
}
