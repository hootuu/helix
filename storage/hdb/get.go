package hdb

import (
	"errors"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Count[T any](dbTx *gorm.DB, query interface{}, args ...interface{}) (int64, error) {
	var m T
	var count int64
	tx := dbTx.Model(&m).Where(query, args...).Count(&count)
	if tx.Error != nil {
		hlog.Err("hdb.Count", zap.Any("query", query), zap.Any("args", args), zap.Error(tx.Error))
		return 0, errors.New("db.count error:" + tx.Error.Error())
	}
	return count, nil
}

func Exist[T any](dbTx *gorm.DB, query interface{}, args ...interface{}) (bool, error) {
	var m T
	var b bool
	tx := dbTx.Model(&m).Select("1").Where(query, args...).Limit(1).Find(&b)
	if tx.Error != nil {
		hlog.Err("hdb.Exist", zap.Any("query", query), zap.Any("args", args), zap.Error(tx.Error))
		return false, errors.New("db.exist error:" + tx.Error.Error())
	}
	return b, nil
}

func ExistWithTable(dbTx *gorm.DB, query interface{}, args ...interface{}) (bool, error) {
	var b bool
	tx := dbTx.Select("1").Where(query, args...).Limit(1).Find(&b)
	if tx.Error != nil {
		hlog.Err("hdb.ExistWithTable", zap.Any("query", query), zap.Any("args", args), zap.Error(tx.Error))
		return false, errors.New("db.exist error:" + tx.Error.Error())
	}
	return b, nil
}

func Get[T any](dbTx *gorm.DB, query string, cond ...interface{}) (*T, error) {
	var m T
	var arr []*T
	tx := dbTx.Model(&m).Where(query, cond...).Limit(1).Find(&arr)
	if tx.Error != nil {
		hlog.Err("hdb.Get", zap.Any("query", query), zap.Any("cond", cond), zap.Error(tx.Error))
		return nil, errors.New("db.get err:" + tx.Error.Error())
	}
	if len(arr) == 0 {
		return nil, nil
	}
	return arr[0], nil
}

func MustGet[T any](dbTx *gorm.DB, query string, cond ...any) (*T, error) {
	obj, err := Get[T](dbTx, query, cond...)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, fmt.Errorf("require: %s %v", query, cond)
	}
	return obj, nil
}
