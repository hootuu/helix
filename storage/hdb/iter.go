package hdb

import (
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Iter[T any](tx func() *gorm.DB, call func(m *T) error) error {
	for {
		arrM, err := Find[T](tx)
		if err != nil {
			hlog.Err("helix.db.iter.find", zap.Error(err))
			return err
		}
		if len(arrM) == 0 {
			break
		}
		for _, sessM := range arrM {
			err := call(sessM)
			if err != nil {
				hlog.Err("helix.db.iter.call", zap.Error(err))
				return err
			}
		}
	}

	return nil
}
