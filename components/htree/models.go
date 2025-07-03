package htree

import (
	"github.com/hootuu/helix/storage/hdb"
)

type TreeM struct {
	hdb.Basic
	ID       ID          `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Sequence int64       `gorm:"column:sequence;"`
	Version  hdb.Version `gorm:"column:version;"`
}
