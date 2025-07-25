package croupier

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hsys"
)

func init() {
	helix.AfterStartup(func() {
		err := zplt.HelixDB().DB().AutoMigrate(&TokenM{})
		if err != nil {
			hsys.Exit(err)
		}
	})
}
