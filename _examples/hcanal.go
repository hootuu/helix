package main

import (
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hcanal"
	"github.com/hootuu/hyle/data/hjson"
)

type TestAlterHandler struct{}

func (t *TestAlterHandler) Schema() string {
	return "helix_mysql"
}

func (t *TestAlterHandler) Table() []string {
	return []string{"harmonic_nineloc_token"}
}

func (t *TestAlterHandler) Action() []string {
	return nil
}

func (t *TestAlterHandler) OnAlter(alter *hcanal.Alter) error {
	fmt.Println("On Alter", hjson.MustToString(alter))
	return nil
}

func main() {
	helix.AfterStartup(func() {
		hcanal.New("canal_main")
	})
	helix.Startup()
}
