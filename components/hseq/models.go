package hseq

import (
	"github.com/hootuu/helix/storage/hdb"
)

type SeqM struct {
	hdb.Basic
	Biz     string      `gorm:"column:biz;primaryKey;size:128;"`
	Seq     int64       `gorm:"column:seq;"`
	Version hdb.Version `gorm:"column:version;"`
}

func (SeqM) TableName() string {
	return "helix_hseq"
}
