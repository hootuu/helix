package croupier

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"time"
)

type TokenM struct {
	hdb.Basic
	ID         ID            `gorm:"column:id;primaryKey;size:32;"`
	Collar     collar.Collar `gorm:"column:collar;not null;size:200;"`
	Bucket     int64         `gorm:"column:bucket;"`
	Remainder  int64         `gorm:"column:remainder;"`
	Expiration *time.Time    `gorm:"column:expiration;"`
}

func (m *TokenM) TableName() string {
	return "helix_croupier_token"
}
