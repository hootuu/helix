package hchan

import (
	"github.com/hootuu/helix/storage/hdb"
	"gorm.io/datatypes"
)

type ChanM struct {
	hdb.Basic
	ID        ID             `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Parent    ID             `gorm:"column:parent;uniqueIndex:uk_parent_name;"`
	Name      string         `gorm:"column:name;uniqueIndex:uk_parent_name;not null;size:60;"`
	Icon      string         `gorm:"column:icon;size:300;"`
	Seq       int            `gorm:"column:seq;"`
	Available bool           `gorm:"column:available"`
	Ctrl      []byte         `gorm:"column:ctrl;size:128;"`
	Tag       datatypes.JSON `gorm:"column:tag;type:jsonb;"`
	Meta      datatypes.JSON `gorm:"column:meta;type:jsonb;"`
}

func (m *ChanM) To() *Channel {
	return &Channel{
		ID:       m.ID,
		Name:     m.Name,
		Icon:     m.Icon,
		Seq:      m.Seq,
		Children: make([]*Channel, 0),
	}
}
