package hnid

import (
	"context"
	"github.com/hootuu/helix/components/hnid/seq"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func init() {
	helix.Use(
		helix.BuildHelix(
			"hnid",
			func() (context.Context, error) {
				err := zplt.HelixPgDB().PG().AutoMigrate(&seq.SequenceM{})
				if err != nil {
					hlog.Err("helix.hnid.init", zap.Error(err))
					return nil, err
				}
				return nil, nil
			},
			func(ctx context.Context) {},
		),
	)
}
