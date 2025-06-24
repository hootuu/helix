package main

import (
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/helix"
	"time"
)

func main() {
	if true {
		factoryTest()
	}
	helix.AfterStartup(func() {
		exampleTree, err := htree.NewTree("example_v5", 5, []uint{2, 4, 4, 5})
		if err != nil {
			fmt.Println(err)
			return
		}
		curId := int64(5030000000000000) //exampleTree.Root()
		s := time.Now()
		count := 100
		for i := 0; i < count; i++ {
			err = exampleTree.Next(curId, func(id htree.ID) error {
				//fmt.Println("nxt: ", id)
				if curId == exampleTree.Root() {
					curId = id
				}
				if i == 2 {
					return errors.New("network err")
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Printf("count: %d, elapse: %d\n", count, time.Now().Sub(s).Milliseconds()/1000)
	})
	helix.Startup()
}

func factoryTest() {
	cfg := []uint{2, 3, 5, 3}
	idGen, _ := htree.NewIdFactory(cfg)
	path, err := idGen.Path(htree.ID(88913312345666))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(path)

	childrenFunc := func(cfg []uint, id int64) {
		idGen, _ := htree.NewIdFactory(cfg)
		fmt.Println("childrenFunc: Deep: ", id, idGen.Deep(id))
		min, max, err := idGen.Children(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("childrenFunc: ", cfg, id, min, max)
	}
	childrenFunc([]uint{2, 3, 5, 3}, 88913300001001)
	childrenFunc([]uint{2, 3, 5, 3}, 88913300001000)
	childrenFunc([]uint{2, 3, 5, 3}, 88913100000000)
	childrenFunc([]uint{2, 3, 5, 3}, 88900000000000)
	childrenFunc([]uint{1, 2, 3, 4}, 83000000000)
	childrenFunc([]uint{2, 2, 3, 3, 3, 4}, 811223334445556666)

	directChildrenFunc := func(cfg []uint, id int64) {
		idGen, _ := htree.NewIdFactory(cfg)
		min, max, base, err := idGen.DirectChildren(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("directChildrenFunc: ", cfg, id, min, max, base)
	}
	directChildrenFunc([]uint{1, 2, 3, 4}, 81223330000)
	directChildrenFunc([]uint{1, 2, 3, 4}, 81220030000)
	directChildrenFunc([]uint{1, 2, 3, 4}, 81220000000)
	directChildrenFunc([]uint{1, 2, 3, 4}, 81000000000)
	directChildrenFunc([]uint{1, 2, 3, 4}, 83000000000)

	idChildrenFunc := func() {
		s := time.Now()
		count := int64(0)
		for i := int64(88900000000000); i < 89000000000000; i += 1000 {
			id := i
			count += 1
			_, _, err := idGen.Children(id)
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println(min, max)
		}
		fmt.Printf("count: %d; seconds: %f\n", count, time.Now().Sub(s).Seconds())
		fmt.Println("elapse: ", time.Now().Sub(s).Milliseconds()/count)
	}
	if false {
		idChildrenFunc()
	}

}
