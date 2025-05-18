package sattva

import (
	"github.com/hootuu/helix/components/sattva/channel"
	"github.com/hootuu/helix/storage/hpg"
	"gorm.io/datatypes"
)

type IdentificationM struct {
	hpg.Basic
	ID   string         `gorm:"column:id;primaryKey;not null;size:32"`
	Info datatypes.JSON `gorm:"type:jsonb;"`
}

func (model *IdentificationM) TableName() string {
	return "helix_sattva_uni_identification"
}

type ChannelM struct {
	hpg.Basic
	ID        channel.ID     `gorm:"column:id;primaryKey;not null;size:32"`
	Type      channel.Type   `gorm:"column:channel_type;"`
	Code      string         `gorm:"column:channel_code;not null;size:50"`
	Config    datatypes.JSON `gorm:"column:config;type:jsonb;"`
	Available bool           `gorm:"column:available;"`
}

func (model *ChannelM) TableName() string {
	return "helix_sattva_uni_channel"
}
