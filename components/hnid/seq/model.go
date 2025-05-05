package seq

import (
	"github.com/hootuu/helix/storage/hpg"
)

type SequenceM struct {
	hpg.Basic
	Code       string      `gorm:"column:code;primaryKey;not null;size:128"`
	MinStart   uint64      `gorm:"column:min_start"`
	MaxEnd     uint64      `gorm:"column:max_end"`
	Step       uint64      `gorm:"column:step"`
	CurrentSeq uint64      `gorm:"column:current_seq"`
	Version    hpg.Version `gorm:"column:version"`
}

func (m *SequenceM) TableName() string {
	return "helix_hnid_seq"
}
