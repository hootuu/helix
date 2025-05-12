package hguard

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func init() {
	helix.Use(helix.BuildHelix(
		"hguard",
		func() (context.Context, error) {
			return nil, zplt.HelixPgDB().PG().AutoMigrate(&GuardM{})
		},
		func(ctx context.Context) {

		},
	))
}
