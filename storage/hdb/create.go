package hdb

import (
	"errors"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Create[T any](dbTx *gorm.DB, model *T) error {
	return doDbActWithRetry(func() error {
		tx := dbTx.Create(model)
		if tx.Error != nil {
			hlog.Err("hdb.Create", zap.Any("model", model), zap.Error(tx.Error))
			return errors.New("db create error: " + tx.Error.Error())
		}
		return nil
	})
}

func MultiCreate[T any](dbTx *gorm.DB, arr []*T) error {
	return doDbActWithRetry(func() error {
		var model T
		tx := dbTx.Model(&model).Create(&arr)
		if tx.Error != nil {
			hlog.Err("hdb.MultiCreate", zap.Any("arr", arr), zap.Error(tx.Error))
			return errors.New("db create multi error: " + tx.Error.Error())
		}
		return nil
	})
}

func GetOrCreate[T any](dbTx *gorm.DB, model *T, cond ...any) error {
	err := dbTx.FirstOrCreate(model, cond...).Error
	if err != nil {
		hlog.Err("hdb.GetOrCreate", zap.Any("model", model), zap.Error(err))
		return err
	}
	return nil
}
