package hlink

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/components/zplt/zcanal"
	"github.com/hootuu/helix/components/zplt/zmeili"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hcanal"
	"github.com/hootuu/helix/storage/hmeili"
)

func init() {
	helix.Use(helix.BuildHelix("helix_link", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&LinkCodeM{},
			&LinkM{},
		)
		if err != nil {
			return nil, err
		}
		meiliPtr := zmeili.HelixMeili()
		indexer := &linkIndexer{}
		err = hmeili.InitIndexer(meiliPtr, indexer)
		if err != nil {
			return nil, err
		}
		zcanal.HelixCanal().RegisterAlterHandler(
			hcanal.NewIndexHandler((&LinkM{}).TableName(), indexer, meiliPtr),
		)
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
