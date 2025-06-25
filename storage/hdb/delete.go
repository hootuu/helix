package hdb

import (
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Delete[T any](dbTx *gorm.DB, cond ...any) error {
	var m T
	tx := dbTx.Unscoped().Delete(&m, cond)
	if tx.Error != nil {
		hlog.Err("hdb.Delete error", zap.Error(tx.Error))
		return tx.Error
	}
	return nil
}
