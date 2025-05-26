package htree

import (
	"github.com/hootuu/helix/storage/hpg"
)

type TreeM struct {
	hpg.Basic
	ID       ID          `gorm:"column:id;primaryKey;"`
	Sequence int64       `gorm:"column:sequence;"`
	Version  hpg.Version `gorm:"column:version;"`
}
