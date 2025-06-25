package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/horg"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		level01ID, err := horg.Create(context.Background(), horg.CreateParas{
			Biz:   collar.Build("biz_test", "biz-"+cast.ToString(time.Now().Unix())),
			Alias: "",
			Name:  "一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十",
			Meta:  nil,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(level01ID)
		level02ID, err := horg.Add(context.Background(), horg.AddParas{
			Parent: level01ID,
			Alias:  "",
			Name:   "一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十",
			Meta:   nil,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(level02ID)
		level03ID, err := horg.Add(context.Background(), horg.AddParas{
			Parent: level02ID,
			Alias:  "",
			Name:   "一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十",
			Meta:   nil,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(level03ID)
		level04ID, err := horg.Add(context.Background(), horg.AddParas{
			Parent: level03ID,
			Alias:  "",
			Name:   "一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十",
			Meta:   nil,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(level04ID)
		//level05ID, err := horg.Add(context.Background(), horg.AddParas{
		//	Parent: level04ID,
		//	Alias:  "",
		//	Name:   "一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十",
		//	Meta:   nil,
		//})
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(level05ID)
		ctx := context.WithValue(context.Background(), "_trace_id_", "hello_world_"+cast.ToString(time.Now().Unix()))
		err = hdb.Tx(zplt.HelixPgCtx(ctx), func(tx *gorm.DB) error {
			var authIDArr []horg.AuthID
			for i := 0; i < 1; i++ {
				authID, err := horg.AddAuth(hdb.TxCtx(tx, ctx), level01ID, "ADMIN-"+cast.ToString(i), []string{"abc.create8888", "abc.delete8888"})
				if err != nil {
					return err
				}
				authIDArr = append(authIDArr, authID)
			}

			err = horg.BindMember(hdb.TxCtx(tx, ctx), level01ID, "ABC-USER", authIDArr)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic(err)
		}

	})
	helix.Startup()
}
