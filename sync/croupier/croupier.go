package croupier

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/hlock"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyle/hypes/collar"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Croupier struct {
	link collar.Collar
	id   ID
}

func Light(uniLink collar.Collar) *Croupier {
	id := BuildID(uniLink)
	return &Croupier{id: id, link: uniLink}
}

func (c *Croupier) Publish(
	ctx context.Context,
	bucket int64,
	ttl time.Duration,
) error {
	tokenM, err := hdb.Get[TokenM](zplt.HelixPgCtx(ctx), "id = ?", c.id)
	if err != nil {
		hlog.TraceErr("croupier.Light: hdb.Get failed", ctx, err,
			zap.String("collar", c.link.ToString()),
			zap.String("id", c.id))
		return err
	}

	bLockDo, err := hlock.Light().Lock(ctx,
		c.getInitLockKey(),
		func() error {
			if tokenM == nil {
				tokenM = &TokenM{
					ID:        c.id,
					Collar:    c.link,
					Bucket:    bucket,
					Remainder: bucket,
				}
				if ttl > 0 {
					expr := time.Now().Add(ttl)
					tokenM.Expiration = &expr
				}
				err := hdb.Create[TokenM](zplt.HelixDB().DB(), tokenM)
				if err != nil {
					return err
				}
				return nil
			}
			mut := map[string]any{
				"bucket":    bucket,
				"remainder": bucket,
			}
			if ttl > 0 {
				expr := time.Now().Add(ttl)
				mut["expiration"] = &expr
			}
			err := hdb.Update[TokenM](zplt.HelixDB().DB(), mut, "id = ?", c.id)
			if err != nil {
				return err
			}
			return nil
		}, 5*time.Minute)
	if err != nil {
		hlog.TraceErr("croupier.Light: OnceLock failed", ctx, err,
			zap.String("collar", c.link.ToString()),
			zap.String("id", c.id))
		return err
	}

	if bLockDo {
		return nil
	}

	err = hretry.Must(func() error {
		tokenM, err = hdb.Get[TokenM](zplt.HelixPgCtx(ctx), "id = ?", c.id)
		if err != nil {
			hlog.TraceErr("croupier.Light[twice]: hdb.Get failed", ctx, err,
				zap.String("collar", c.link.ToString()),
				zap.String("id", c.id))
			return err
		}
		if tokenM != nil {
			return nil
		}
		return fmt.Errorf("init helix.croupier failed")
	})
	if err != nil {
		hlog.TraceErr("croupier.Light[twice]: hdb.Get failed", ctx, err,
			zap.String("collar", c.link.ToString()),
			zap.String("id", c.id))
		return err
	}
	return nil
}

func (c *Croupier) ID() string {
	return c.id
}

func (c *Croupier) Link() collar.Collar {
	return c.link
}

func (c *Croupier) Allow(ctx context.Context, call func() error) (allow bool, err error) {
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx,
			fmt.Sprintf("helix.croupier[%s]", c.link),
			hlog.F(zap.String("id", c.id)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return []zap.Field{zap.Bool("allow", allow)}
			},
		)()
	}

	bLockDo, err := hlock.Light().Lock(ctx,
		c.getEachLockKey(),
		func() error {
			allow, err = c.doAllowInLock(call)
			if err != nil {
				return err
			}
			return nil
		}, 10*time.Second)
	if err != nil {
		return false, err
	}
	if !bLockDo {
		return false, nil
	}

	return allow, nil
}

func (c *Croupier) doAllowInLock(call func() error) (bool, error) {
	rows, err := hdb.UpdateX[TokenM](
		zplt.HelixDB().DB(),
		map[string]any{
			"remainder": gorm.Expr("remainder - 1"),
		},
		"id = ? AND remainder > 0",
		c.id,
	)
	if err != nil {
		return false, err
	}
	if rows == 0 {
		return false, nil
	}
	err = call()
	if err != nil {
		hlog.TraceErr("croupier.doAllowInLock: call failed", nil, err)
		//reback the bucket
		innerErr := hdb.Update[TokenM](
			zplt.HelixDB().DB(),
			map[string]any{
				"remainder": gorm.Expr("remainder + 1"),
			},
			"id = ? AND remainder < ?",
			c.id, gorm.Expr("bucket"),
		)
		if innerErr != nil {
			return false, innerErr
		}
		return false, err
	}
	return true, nil
}

func (c *Croupier) getInitLockKey() string {
	return fmt.Sprintf("helix_croupier:init:%s", c.id)
}
func (c *Croupier) getEachLockKey() string {
	return fmt.Sprintf("helix_croupier:each:%s", c.id)
}
