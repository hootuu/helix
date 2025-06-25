package hseq

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
)

type SeqM struct {
	hdb.Basic
	Biz     collar.ID   `gorm:"column:biz;primaryKey;"`
	Seq     int64       `gorm:"column:seq;"`
	Version hdb.Version `gorm:"column:version;"`
}

func (SeqM) TableName() string {
	return "helix_hseq"
}
