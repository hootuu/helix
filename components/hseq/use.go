package hseq

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/hlock"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/gorm"
	"time"
)

var gLocker *hlock.Locker

var errVersionConflict = errors.New("optimistic lock version conflict")

func Next(ctx context.Context, biz collar.Collar) (ID, error) {
	bizID := biz.ToSafeID()
	var nxtID ID

	const maxRetries = 3
	const retryInterval = 50 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		bLockDo, err := gLocker.Lock(ctx, "hseq:"+bizID, func() error {
			tx := zplt.HelixPgCtx(ctx)
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
				"seq":     gorm.Expr("seq + 1"),
				"version": gorm.Expr("version + 1"),
			}
			rows, err := hdb.UpdateX[SeqM](tx, mut, "biz = ? AND version = ?", bizID, seqM.Version)
			if err != nil {
				return err
			}
			if rows == 0 {
				return errVersionConflict
			}

			nxtID = seqM.Seq + 1
			return nil
		}, 3*time.Second)

		if err != nil {
			if errors.Is(err, errVersionConflict) {
				continue
			}
			return 0, fmt.Errorf("unexpected error during lock execution for biz %s: %w", bizID, err)
		}

		if bLockDo {
			return nxtID, nil // 成功获取ID，直接返回
		}
		time.Sleep(retryInterval)
	}
	return 0, fmt.Errorf("failed to get next ID for biz %s after %d retries", bizID, maxRetries)
}

func Current(ctx context.Context, biz collar.Collar, call func(current ID)) error {
	bizID := biz.ToSafeID()
	tx := zplt.HelixPgCtx(ctx)
	bLockDo, err := gLocker.Lock(ctx, "hseq:"+bizID, func() error {
		seqM, err := hdb.Get[SeqM](tx, "biz = ?", bizID)
		if err != nil {
			return err
		}
		if seqM == nil {
			call(0)
		} else {
			call(seqM.Seq)
		}
		return nil
	}, 200*time.Millisecond)
	if err != nil {
		return err
	}
	if !bLockDo {
		return errors.New("biz " + bizID + ": lock failed")
	}
	return nil
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
