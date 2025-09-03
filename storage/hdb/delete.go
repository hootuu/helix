package hdb

import (
	"errors"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Delete[T any](dbTx *gorm.DB, cond ...any) error {
	var m T
	tx := dbTx.Unscoped().Delete(&m, cond...)
	if tx.Error != nil {
		hlog.Err("hdb.Delete error", zap.Error(tx.Error))
		return tx.Error
	}
	return nil
}

func DeleteLogic[T any](dbTx *gorm.DB, cond ...any) error {
	var m T
	if len(cond) == 0 {
		return errors.New("SoftDelete called without condition")
	}

	tx := dbTx.Delete(&m, cond...)
	if tx.Error != nil {
		hlog.Err("hdb.SoftDelete error", zap.Error(tx.Error))
		return tx.Error
	}
	return nil
}
