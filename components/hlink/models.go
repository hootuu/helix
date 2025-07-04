package hlink

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
)

type LinkCodeM struct {
	hdb.Basic
	Link collar.ID `gorm:"column:link;primaryKey;not null;size:64;"`
	Biz  string    `gorm:"column:biz;uniqueIndex:uk_biz_code;not null;size:64;"`
	Code string    `gorm:"column:code;uniqueIndex:uk_biz_code;not null;size:64;"`
}

func (m *LinkCodeM) TableName() string {
	return "helix_link_code"
}
