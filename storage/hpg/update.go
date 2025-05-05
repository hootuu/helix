package hpg

import (
	"errors"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Update[T any](dbTx *gorm.DB, mut map[string]any, query any, args ...any) error {
	var m T
	tx := dbTx.Model(&m).Where(query, args...).Updates(mut)
	if tx.Error != nil {
		hlog.Err("hpg.Update", zap.Any("mut", mut),
			zap.Any("query", query), zap.Any("args", args), zap.Error(tx.Error))
		return errors.New("db.update error:" + tx.Error.Error())
	}
	return nil
}
