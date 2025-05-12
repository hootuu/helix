package sattva

import (
	"github.com/hootuu/helix/storage/hpg"
	"gorm.io/datatypes"
	"strings"
)

type IdentityM struct {
	hpg.Basic
	ID       string `gorm:"column:id;primaryKey;not null;size:32"`
	Nickname string `gorm:"column:nickname;not null;size:32"`
	Avatar   string `gorm:"column:avatar;not null;size:200"`
}

func buildIdentityTableName(code string) string {
	return "helix_hwt_" + strings.ToLower(code) + "_identity"
}

type AuthenticatorM struct {
	hpg.Basic
	ID       string         `gorm:"column:id;uniqueIndex:uk_id_type;not null;size:32"`
	AuthType AuthType       `gorm:"column:auth_type;uniqueIndex:uk_id_type;not null;size:16"`
	Data     datatypes.JSON `gorm:"type:jsonb;"`
}

func buildAuthenticatorTableName(code string) string {
	return "helix_hwt_" + strings.ToLower(code) + "_authenticator"
}
