package hlink

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func init() {
	helix.Use(helix.BuildHelix("helix_link", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&LinkCodeM{},
		)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
