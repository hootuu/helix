package hvault

import "github.com/hootuu/helix/storage/hdb"

type VaultM struct {
	hdb.Basic
	Idx        string `gorm:"column:idx;primaryKey;not null;size:32"`
	PrivateKey []byte `gorm:"column:private_key"`
	Usage      uint64 `gorm:"column:usage"`
	Available  bool   `gorm:"column:available"`
}

func (m *VaultM) TableName() string {
	return "helix_hvault"
}
