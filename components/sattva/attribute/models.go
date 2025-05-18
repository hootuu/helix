package attribute

import (
	"github.com/hootuu/helix/storage/hpg"
	"gorm.io/datatypes"
)

type BasicM struct {
	hpg.Basic
	Identification string `gorm:"column:identification;uniqueIndex:uk_id_attr;not null;size:32"`
	Attr           string `gorm:"column:attr;uniqueIndex:uk_id_attr;not null;size:64"`
}

type SimpleM struct {
	BasicM
	Value string `gorm:"column:value;not null;size:500"`
}

func (model *SimpleM) TableName() string {
	return "helix_sattva_attribute_simple"
}

type ComplexM struct {
	BasicM
	Value datatypes.JSON `gorm:"column:value;type:jsonb;"`
}

func (model *ComplexM) TableName() string {
	return "helix_sattva_attribute_complex"
}
