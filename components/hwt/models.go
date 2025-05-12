package hwt

import (
	"github.com/hootuu/helix/storage/hpg"
	"strings"
)

type RefreshTokenM struct {
	hpg.Basic
	RefreshToken string `gorm:"column:refresh_token;primaryKey;not null;size:64"`
	Identity     string `gorm:"column:identity;uniqueIndex:uk_id_code;not null;size:32"`
	Code         string `gorm:"column:code;uniqueIndex:uk_id_code;not null;size:32"`
	Expiration   int64  `gorm:"column:expiration;"`
}

func buildRefreshTokenTableName(code string) string {
	return "helix_hwt_" + strings.ToLower(code) + "_refresh_token"
}

//
//type RefreshTokenM struct {
//	hpg.Basic
//	GuardID        string    `gorm:"column:biz;index;not null;size:32"`
//	Token          string    `gorm:"column:id;primaryKey;not null;size:64"`
//	TokenTimestamp time.Time `gorm:"column:token_timestamp;"`
//	Usage          int64     `gorm:"column:usage;"`
//	LstUsageTime   time.Time `gorm:"column:lst_usage_time;"`
//}

//func (model *RefreshTokenM) TableName() string {
//	return "helix_guard_refresh_token"
//}

//func buildToken() (string, time.Time) {
//	timestamp := time.Now()
//	_, pri, err := hed25519.Random()
//	if err != nil {
//		hlog.Err("hguard.buildToken", zap.Error(err))
//		return idx.New(), timestamp
//	}
//	return base58.Encode(pri), timestamp
//}
