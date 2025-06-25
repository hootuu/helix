package channel

import (
	"github.com/hootuu/helix/storage/hdb"
	"gorm.io/datatypes"
)

type IdChannelM struct {
	hdb.Basic
	ID      string         `gorm:"column:id;index;not null;size:32"`
	Channel ID             `gorm:"column:channel;uniqueIndex:uk_channel_link;not null;size:32"`
	Link    Link           `gorm:"column:link;uniqueIndex:uk_channel_link;not null;size:128"`
	Paras   datatypes.JSON `gorm:"column:paras;type:jsonb;"`
}

func (model *IdChannelM) TableName() string {
	return "helix_sattva_uni_id_channel"
}
