package main

import (
	"fmt"
	"github.com/hootuu/helix/components/hchan"
	"github.com/hootuu/helix/helix"
	"github.com/spf13/cast"
)

func main() {
	helix.AfterStartup(func() {
		f, err := hchan.NewFactory("test")
		if err != nil {
			panic(err)
		}
		for i := 0; i < 100; i++ {
			c, err := f.Add(0, "CHANNEL-"+cast.ToString(i), "https://icon.cn", i)
			if err != nil {
				panic(err)
			}
			fmt.Println(c)
			for j := 0; j < 100; j++ {
				c2, err := f.Add(c, "CHANNEL-"+cast.ToString(i)+"-"+cast.ToString(j), "https://icon.cn", j)
				if err != nil {
					panic(err)
				}
				fmt.Println(c2)
			}
		}

	})
	helix.Startup()
}
