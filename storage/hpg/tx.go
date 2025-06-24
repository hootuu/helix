package hpg

import (
	"context"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	hpgTxInCtx = "_hpg_tx_"
)

func TxCtx(tx *gorm.DB, ctx ...context.Context) context.Context {
	var parent context.Context
	if len(ctx) > 0 {
		parent = ctx[0]
	} else {
		parent = context.Background()
	}
	return context.WithValue(parent, hpgTxInCtx, tx)
}

func CtxTx(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return nil
	}
	obj := ctx.Value(hpgTxInCtx)
	if obj == nil {
		return nil
	}
	return obj.(*gorm.DB)
}

func Tx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	nErr := db.Transaction(func(tx *gorm.DB) error {
		err := fn(tx)
		if err != nil {
			hlog.Err("hpg.Tx: fn", zap.Error(err))
			return err
		}
		return nil
	})
	if nErr != nil {
		hlog.Err("hpg.Tx: db.Transaction", zap.Error(nErr))
		return nErr
	}
	return nil
}
