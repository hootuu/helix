package hpg

import (
	"errors"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Exist[T any](dbTx *gorm.DB, query interface{}, args ...interface{}) (bool, error) {
	var m T
	var b bool
	tx := dbTx.Model(&m).Select("1").Where(query, args...).Limit(1).Find(&b)
	if tx.Error != nil {
		hlog.Err("hpg.Exist", zap.Any("query", query), zap.Any("args", args), zap.Error(tx.Error))
		return false, errors.New("db.exist error:" + tx.Error.Error())
	}
	return b, nil
}

func Get[T any](dbTx *gorm.DB, query string, cond ...interface{}) (*T, error) {
	var m T
	var arr []*T
	tx := dbTx.Model(&m).Where(query, cond...).Limit(1).Find(&arr)
	if tx.Error != nil {
		hlog.Err("hpg.Get", zap.Any("query", query), zap.Any("cond", cond), zap.Error(tx.Error))
		return nil, errors.New("db.get err:" + tx.Error.Error())
	}
	if len(arr) == 0 {
		return nil, nil
	}
	return arr[0], nil
}
