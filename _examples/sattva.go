package main

import (
	"fmt"
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hcast"
	"time"
)

func main() {

	helix.AfterStartup(func() {
		const bizCode = "pwd.example.biz"
		s, err := sattva.NewSattva("example_sattva")
		if err != nil {
			fmt.Println(err)
			return
		}
		chnID, err := s.RegisterChannel(sattva.Password,
			bizCode,
			dict.NewDict().Set("encrypt_password", "88888888"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("chnID: ", chnID)
		total := 1
		start := time.Now()
		for i := 0; i < total; i++ {
			_, err := s.IdentificationCreate(&sattva.Channel{
				Type:  sattva.Password,
				Code:  bizCode,
				Link:  fmt.Sprintf("user_%d_%d", i, time.Now().UnixMilli()),
				Paras: dict.New(make(map[string]interface{})).Set("password", "12345678"),
			}, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println("id: ", id)
		}
		fmt.Println("elapse: ", int64(time.Now().Sub(start))/int64(time.Millisecond)/int64(total))
		ms := time.Now().UnixMilli()
		id, err := s.IdentificationCreate(&sattva.Channel{
			Type:  sattva.Password,
			Code:  bizCode,
			Link:  fmt.Sprintf("user_example_%d", ms),
			Paras: dict.New(make(map[string]interface{})).Set("password", "12345678"),
		}, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("id: ", id)
		willOK, id, err := s.Identify(&sattva.Channel{
			Type:  sattva.Password,
			Code:  bizCode,
			Link:  fmt.Sprintf("user_example_%d", ms),
			Paras: dict.New(make(map[string]interface{})).Set("password", "12345678"),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("willOK: ", willOK)

		willFail, id, err := s.Identify(&sattva.Channel{
			Type:  sattva.Password,
			Code:  bizCode,
			Link:  fmt.Sprintf("user_example_%d", ms),
			Paras: dict.New(make(map[string]interface{})).Set("password", "abc"),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("willFail: ", willFail)
		attrS := time.Now()
		var nameAttr string
		var ageAttr string
		var lottoAttr string
		for i := 0; i < 10000; i++ {
			if i < 10 {
				nameAttr = "name-" + hcast.ToString(time.Now().UnixMilli()) + "-" + hcast.ToString(i)
				err = s.SetAttribute(id, nameAttr, "NAME-TEST")
				if err != nil {
					fmt.Println(err)
					return
				}

				ageAttr = "age-" + hcast.ToString(time.Now().UnixMilli()) + "-" + hcast.ToString(i)
				err = s.SetAttribute(id, ageAttr, 18)
				if err != nil {
					fmt.Println(err)
					return
				}
				lottoAttr = "lotto-" + hcast.ToString(time.Now().UnixMilli()) + "-" + hcast.ToString(i)
				err = s.SetAttribute(id, lottoAttr, map[string]interface{}(dict.NewDict().Set("ratio", 1090).
					Set("token", "CNY").Set("lotto", &struct {
					Day       int
					Week      int
					Timestamp int64
				}{
					Day:       99,
					Week:      9,
					Timestamp: time.Now().UnixMilli(),
				})))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			_, err := s.GetAttribute(id, nameAttr, ageAttr, lottoAttr)
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println(hjson.MustToString(info))
		}
		fmt.Println("elapse:", (time.Now().UnixMilli()-attrS.UnixMilli())/1000)

	})
	helix.Startup()
}
