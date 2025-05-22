package main

import (
	"fmt"
	"github.com/hootuu/helix/components/htree"
)

func main() {
	idGen, _ := htree.NewIdFactory([]uint{2, 3, 5})
	path, err := idGen.Path(htree.ID(88913312345))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(path)
}
