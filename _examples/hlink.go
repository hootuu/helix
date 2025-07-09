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
		s := time.Now()
		for i := 0; i < 10000; i++ {

			majorID := "uid_0000009_" + cast.ToString(time.Now().UnixNano())
			major := collar.Build("USER", majorID)
			ouser := collar.Build("USER", "uid_000000222225555")
			ctx := context.Background()
			inviteCode, err := hlink.Generate(
				ctx,
				"NINEORA_USER_INVITE",
				major,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println("invite code:", inviteCode)

			err = hlink.Bind(
				ctx,
				"NINEORA_USER_INVITE",
				inviteCode,
				"FRIEND",
				ouser,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			//
			//time.Sleep(800 * time.Millisecond)
			//
			//fmt.Println("majorID: ", majorID)
			//p, err := hlink.Filter(
			//	"major_code = 'USER' AND major_id = '"+majorID+"'",
			//	[]string{"auto_id:desc"},
			//	pagination.PageNormal(),
			//)
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//fmt.Println(hjson.MustToString(p))
			//
			//time.Sleep(1 * time.Second)
			//err = hlink.Unbind(
			//	ctx,
			//	"NINEORA_USER_INVITE",
			//	major,
			//	"FRIEND",
			//	ouser,
			//)
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//fmt.Println("unbind ok")
			//
			//err = hlink.Bind(
			//	ctx,
			//	"NINEORA_USER_INVITE",
			//	inviteCode,
			//	"FRIEND",
			//	ouser,
			//)
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//
			//fmt.Println("bind ok")
			//
			//err = hlink.Bind(
			//	ctx,
			//	"NINEORA_USER_INVITE",
			//	inviteCode,
			//	"FRIEND",
			//	ouser,
			//)
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//
			//fmt.Println("bind err")

		}

		fmt.Println(time.Since(s).Milliseconds() / 10000)
	})

	helix.Startup()
}
