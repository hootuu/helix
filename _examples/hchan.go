package main

import (
	"fmt"
	"github.com/hootuu/helix/components/hchan"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		hchanExample()
		hchanRunning()
	})
	helix.Startup()
}

func hchanExample() {
	f, err := hchan.NewFactory("test")
	if err != nil {
		panic(err)
	}
	cID, err := f.Add(0, "CHANNEL-"+cast.ToString(time.Now().Unix()), "https://icon.cn", 1)
	if err != nil {
		panic(err)
	}
	err = f.Mut(cID,
		"NEW-CHANNEL-"+cast.ToString(time.Now().Unix()),
		"https://new.ic",
		3444,
	)
	if err != nil {
		panic(err)
	}

	err = f.SetAvailable(cID, false)
	if err != nil {
		panic(err)
	}

	err = f.SetAvailable(cID, true)
	if err != nil {
		panic(err)
	}

	err = f.SetSeq(cID, 1001)
	if err != nil {
		panic(err)
	}

	ch, err := f.GetChildren(0, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(hjson.MustToString(ch))

	time.Sleep(100 * time.Millisecond)
	pageData, err := f.Filter(
		"parent = '"+cast.ToString(f.Root())+"'",
		[]string{},
		pagination.NewPage(1, 1),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(hjson.MustToString(pageData))
}

func hchanRunning() {
	f, err := hchan.NewFactory("test2")
	if err != nil {
		panic(err)
	}
	s := time.Now()
	for i := 0; i < 100; i++ {
		c, err := f.Add(0, "CHANNEL-"+cast.ToString(i), "https://icon.cn", i)
		if err != nil {
			panic(err)
		}
		//fmt.Println(c)
		for j := 0; j < 100; j++ {
			_, err := f.Add(c, "CHANNEL-"+cast.ToString(i)+"-"+cast.ToString(j), "https://icon.cn", j)
			if err != nil {
				panic(err)
			}
			//fmt.Println(c2)
		}
	}
	e := time.Since(s).Milliseconds()
	fmt.Println(e, " ->", e/10000)
}
