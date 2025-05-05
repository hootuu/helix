package hidem

import (
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
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
		table:    f.tableName(),
	}
	exist, err := hpg.Exist[IdemM](zplt.HelixPgDB().PG(), "idem_code = ?", idemM.IdemCode)
	if err != nil {
		return false, err
	}
	if exist {
		return false, nil
	}
	err = hpg.Create[IdemM](zplt.HelixPgDB().PG(), idemM)
	if err != nil {
		return false, err
	}
	return true, nil
}

func newDbFactory(code string, expiration time.Duration, cleanInterval time.Duration) (*dbFactory, error) {
	f := &dbFactory{
		code:          code,
		expiration:    expiration,
		cleanInterval: cleanInterval,
	}
	err := zplt.HelixPgDB().PG().AutoMigrate(&IdemM{
		table: f.tableName(),
	})
	if err != nil {
		hlog.Err("hidem.pg.newDbFactory", zap.Error(err))
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

		tx := zplt.HelixPgDB().PG().Unscoped().
			Where("created_at < ?", threshold).
			Delete(&IdemM{
				table: f.tableName(),
			})
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
	hpg.Basic
	IdemCode string `gorm:"column:idem_code;primaryKey;not null;size:128"`
	table    string `gorm:"-"`
}

func (m *IdemM) TableName() string {
	return m.table
}
