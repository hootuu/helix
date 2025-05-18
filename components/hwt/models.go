package hwt

import (
	"github.com/hootuu/helix/storage/hpg"
	"strings"
)

type RefreshTokenM struct {
	hpg.Basic
	RefreshToken string `gorm:"column:refresh_token;primaryKey;not null;size:88"`
	Identity     string `gorm:"column:identity;uniqueIndex:uk_id_code;not null;size:32"`
	Code         string `gorm:"column:code;uniqueIndex:uk_id_code;not null;size:32"`
	Expiration   int64  `gorm:"column:expiration;"`
}

func buildRefreshTokenTableName(code string) string {
	return "helix_hwt_" + strings.ToLower(code) + "_refresh_token"
}
