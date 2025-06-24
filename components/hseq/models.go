package hseq

import (
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hypes/collar"
)

type SeqM struct {
	hpg.Basic
	Biz     collar.ID   `gorm:"column:biz;primaryKey;"`
	Seq     int64       `gorm:"column:seq;"`
	Version hpg.Version `gorm:"column:version;"`
}

func (SeqM) TableName() string {
	return "helix_hseq"
}
