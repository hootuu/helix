package hguard

import (
	"github.com/hootuu/helix/storage/hdb"
	"time"
)

type GuardM struct {
	hdb.Basic
	ID           string    `gorm:"column:id;primaryKey;not null;size:32"`
	Biz          string    `gorm:"column:biz;uniqueIndex:uk_biz_alias;not null;size:32"`
	Alias        string    `gorm:"column:alias;uniqueIndex:uk_biz_alias;not null;size:32"`
	PubKey       []byte    `gorm:"column:pub_key;"`
	PriKey       []byte    `gorm:"column:pri_key;"`
	Usage        int64     `gorm:"column:usage;"`
	LstUsageTime time.Time `gorm:"column:lst_usage_time;"`
}

func (model *GuardM) TableName() string {
	return "helix_guard_guard"
}
