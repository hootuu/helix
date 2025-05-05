package hpg

import (
	"gorm.io/gorm"
	"time"
)

type Version uint64

func (v Version) Inc() Version {
	return v + 1
}

type Basic struct {
	AutoID    int64          `gorm:"column:auto_id;uniqueIndex;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;index;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
