package hdb

import (
	"errors"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Update[T any](dbTx *gorm.DB, mut map[string]any, query any, args ...any) error {
	_, err := UpdateX[T](dbTx, mut, query, args...)
	return err
}

func UpdateX[T any](dbTx *gorm.DB, mut map[string]any, query any, args ...any) (int64, error) {
	var m T
	tx := dbTx.Model(&m).Where(query, args...).Updates(mut)
	if tx.Error != nil {
		hlog.Err("hdb.UpdateX", zap.Any("mut", mut),
			zap.Any("query", query), zap.Any("args", args), zap.Error(tx.Error))
		return 0, errors.New("db.update error:" + tx.Error.Error())
	}
	return tx.RowsAffected, nil
}
