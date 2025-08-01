package hidem

import (
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hync"
	"go.uber.org/zap"
	"time"
)

type dbFactory struct {
	code          string
	expiration    time.Duration
	cleanInterval time.Duration
	syncLine      hync.Line
	lstCleanTime  time.Time
	db            *hdb.Database
}

func (f *dbFactory) Check(idemCode string) (bool, error) {
	if err := CheckIdemCode(idemCode); err != nil {
		return false, err
	}
	defer func() {
		if f.expiration == NoExpiration {
			return
		}
		go f.clean()
	}()
	idemM := &IdemM{
		IdemCode: idemCode,
	}
	exist, err := hdb.Exist[IdemM](f.db.DB().Table(f.tableName()), "idem_code = ?", idemM.IdemCode)
	if err != nil {
		return false, err
	}
	if exist {
		return false, nil
	}
	err = hdb.Create[IdemM](f.db.DB().Table(f.tableName()), idemM)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *dbFactory) MustCheck(idemCode string) error {
	ok, err := f.Check(idemCode)
	if err != nil {
		return err
	}
	if !ok {
		// todo: return error
		hlog.Fix("tx.Mint: idem " + idemCode)
	}
	return nil
}

func newDbFactory(db *hdb.Database, code string, expiration time.Duration, cleanInterval time.Duration) (*dbFactory, error) {
	f := &dbFactory{
		code:          code,
		expiration:    expiration,
		cleanInterval: cleanInterval,
		db:            db,
	}
	err := f.db.DB().Table(f.tableName()).AutoMigrate(&IdemM{})
	if err != nil {
		hlog.Err("hidem.db.newDbFactory", zap.Error(err))
		return nil, err
	}
	return f, nil
}

func (f *dbFactory) tableName() string {
	return fmt.Sprintf("helix_idem_%s", f.code)
}

func (f *dbFactory) clean() {
	_ = f.syncLine.Do(func() error {
		if time.Now().Sub(f.lstCleanTime) < f.cleanInterval {
			return nil
		}
		threshold := time.Now().Add(-f.expiration)
		cleanCount := int64(0)
		var err error

		defer hlog.Elapse("helix.idem.pg.clean", func() []zap.Field {
			return []zap.Field{
				zap.Time("lstCleanTime", f.lstCleanTime),
				zap.Time("threshold", threshold),
			}
		}, func() []zap.Field {
			return []zap.Field{
				zap.Int64("cleanCount", cleanCount),
				zap.Time("lstCleanTime", f.lstCleanTime),
				zap.Error(err),
			}
		})()

		tx := f.db.DB().Unscoped().
			Table(f.tableName()).
			Where("created_at < ?", threshold).
			Delete(&IdemM{})
		if tx.Error != nil {
			hlog.Err("hidem.pg.clean[ignore]", zap.Error(tx.Error))
			err = tx.Error
			return tx.Error
		}
		cleanCount = tx.RowsAffected
		f.lstCleanTime = time.Now()
		return nil
	})
}

type IdemM struct {
	hdb.Basic
	IdemCode string `gorm:"column:idem_code;primaryKey;not null;size:128"`
}
