package hseq

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
)

func Next(ctx context.Context, biz collar.Collar) (ID, error) {
	bizID := biz.ToSafeID()
	tx := zplt.HelixPgCtx(ctx)
	seqM, err := hdb.Get[SeqM](tx, "biz = ?", bizID)
	if err != nil {
		return 0, err
	}
	if seqM == nil {
		seqM = &SeqM{
			Biz:     bizID,
			Seq:     1,
			Version: 1,
		}
		err = hdb.Create[SeqM](tx, seqM)
		if err != nil {
			return 0, err
		}
		return seqM.Seq, nil
	}
	mut := map[string]any{
		"seq":     seqM.Seq + 1,
		"version": seqM.Version + 1,
	}
	err = hdb.Update[SeqM](tx, mut, "biz = ? AND version = ?", bizID, seqM.Version)
	if err != nil {
		return 0, err
	}
	return seqM.Seq + 1, nil
}

func init() {
	helix.Use(helix.BuildHelix("helix_seq", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&SeqM{},
		)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
