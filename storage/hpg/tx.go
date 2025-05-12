package hpg

import (
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
