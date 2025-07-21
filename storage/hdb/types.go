package hdb

import (
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/ex"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Transaction = *gorm.DB

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

type Template struct {
	AutoID    int64          `gorm:"column:auto_id;uniqueIndex;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;index;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Ctrl      []byte         `gorm:"column:ctrl;size:128;"`
	Tag       datatypes.JSON `gorm:"column:tag;type:jsonb;"`
	Meta      datatypes.JSON `gorm:"column:meta;type:jsonb;"`
}

func TemplateFromEx(exM *ex.Ex) Template {
	if exM == nil {
		return Template{}
	}
	return Template{
		Ctrl: exM.Ctrl,
		Tag:  hjson.MustToBytes(exM.Tag),
		Meta: hjson.MustToBytes(exM.Meta),
	}
}
