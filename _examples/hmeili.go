package main

import (
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/data/hjson"
	"time"
)

type MeiliDocTest struct {
	hmeili.DocBasic
}

func (m MeiliDocTest) IndexName() string {
	return "meili_test"
}

func (m MeiliDocTest) GetAutoID() uint64 {
	return 0
}

func (m MeiliDocTest) GetTimestamp() int64 {
	return 0
}

func main() {
	helix.AfterStartup(func() {
		meili := zplt.HelixMeili()
		//_, err := meili.Meili().Index("meili_test").UpdateSortableAttributes(&[]string{"auto_id"})
		//if err != nil {
		//	log.Fatalf("更新可排序字段失败: %v", err)
		//}
		//_, err = meili.Meili().Index("meili_test").UpdateFilterableAttributes(&[]string{"auto_id"})
		//if err != nil {
		//	log.Fatalf("更新可过滤字段失败: %v", err)
		//}
		for i := 0; i < 100; i++ {
			var arr []*MeiliDocTest
			for j := 0; j < 1000; j++ {
				arr = append(arr, &MeiliDocTest{DocBasic: hmeili.DocBasic{
					AutoID:    uint64(time.Now().UnixNano()),
					Timestamp: 101,
				}})
			}
			_, err := meili.Meili().Index("meili_test").AddDocuments(arr, "auto_id")
			if err != nil {
				panic(err)
			}
			//_, err = meili.Meili().WaitForTask(t.TaskUID, time.Millisecond*10)
			//if err != nil {
			//	log.Fatalf("等待任务完成失败: %v", err)
			//}
		}

		p, err := hmeili.Find(meili, "game", &hmeili.SearchRequest{
			//Query: "1751567815699316000",
			Filter: []string{"auto_id = 1751567815699316000"},
		}, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(hjson.MustToString(p))
	})
	helix.Startup()
}
