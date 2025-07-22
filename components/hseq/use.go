package hseq

import (
	"context"
	"errors"
	"github.com/hootuu/helix/components/hlock"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"time"
)

var gLocker *hlock.Locker

func Next(ctx context.Context, biz collar.Collar) (ID, error) {
	bizID := biz.ToSafeID()
	tx := zplt.HelixPgCtx(ctx)
	var nxtID ID
	bLockDo, err := gLocker.Lock(ctx, "hseq:"+bizID, func() error {
		seqM, err := hdb.Get[SeqM](tx, "biz = ?", bizID)
		if err != nil {
			return err
		}
		if seqM == nil {
			seqM = &SeqM{
				Biz:     bizID,
				Seq:     1,
				Version: 1,
			}
			err = hdb.Create[SeqM](tx, seqM)
			if err != nil {
				return err
			}
			nxtID = seqM.Seq
			return nil
		}
		mut := map[string]any{
			"seq":     seqM.Seq + 1,
			"version": seqM.Version + 1,
		}
		rows, err := hdb.UpdateX[SeqM](tx, mut, "biz = ? AND version = ?", bizID, seqM.Version)
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("biz " + bizID + ": has been changed")
		}
		nxtID = seqM.Seq + 1
		return nil
	}, 200*time.Millisecond)
	if err != nil {
		return 0, err
	}
	if !bLockDo {
		return 0, errors.New("biz " + bizID + ": lock failed")
	}
	return nxtID, nil
}

func init() {
	helix.Use(helix.BuildHelix("helix_seq", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&SeqM{},
		)
		if err != nil {
			return nil, err
		}
		gLocker = hlock.NewLocker(zplt.HelixRdsCache())
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
