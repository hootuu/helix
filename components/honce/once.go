package honce

import (
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Status int

const (
	EXECUTING Status = 1
	FAILED    Status = -1
	SUCCESS   Status = 8
)

type OnceM struct {
	hpg.Basic
	OnceCode    string      `gorm:"column:once_code;primaryKey;not null;size:128"`
	DoServerID  string      `gorm:"column:do_serv_id;index;not null;size:32"`
	DoStatus    Status      `gorm:"column:do_status"`
	DoStartTime time.Time   `gorm:"column:do_start_time;;not null"`
	DoEndTime   time.Time   `gorm:"column:do_end_time;;not null"`
	Version     hpg.Version `gorm:"column:version;default:0"`
}

func (m *OnceM) TableName() string {
	return "helix_honce"
}

func doSetEnd(m *OnceM, status Status) error {
	return hpg.Update[OnceM](
		zplt.HelixPgDB().PG(),
		map[string]any{
			"do_status":   status,
			"do_end_time": gorm.Expr("CURRENT_TIMESTAMP"),
			"version":     m.Version.Inc(),
		},
		"once_code = ? AND version = ?",
		m.OnceCode, m.Version,
	)
}

func doOnce(onceCode string, call func() error) error {
	onceM, err := hpg.Get[OnceM](zplt.HelixPgDB().PG(), "once_code = ?", onceCode)
	if err != nil {
		hlog.Err("helix.once.doOnce: Get", zap.String("code", onceCode), zap.Error(err))
		return err
	}
	if onceM != nil {
		if onceM.DoStatus == FAILED {
			hlog.Err("helix.once.doOnce: dbOnce==FAILED", zap.String("code", onceCode))
			return errors.New(fmt.Sprintf("a single task failed to execute on other machines: %s", onceCode))
		}
		return nil
	}

	onceM = &OnceM{
		OnceCode:    onceCode,
		DoServerID:  hsys.ServerID(),
		DoStatus:    EXECUTING,
		DoStartTime: time.Now(),
	}

	err = hpg.Create[OnceM](zplt.HelixPgDB().PG(), onceM)
	if err != nil {
		return err
	}

	err = call()
	if err != nil {
		hlog.Err("helix.once.doOnce: local execute failed", zap.String("code", onceCode), zap.Error(err))
		igErr := doSetEnd(onceM, FAILED)
		if igErr != nil {
			hlog.Err("helix.once.doOnce:set status(failed) failed",
				zap.String("code", onceCode), zap.Error(igErr))
		}
		return err
	}

	igErr := doSetEnd(onceM, SUCCESS)
	if igErr != nil {
		hlog.Err("helix.once.doOnce:set status(success) failed", zap.String("code", onceCode), zap.Error(igErr))
	}

	return nil
}
