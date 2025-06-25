package hcaptcha

import (
	"github.com/hootuu/helix/storage/hdb"
	"time"
)

type CaptchaM struct {
	hdb.Basic
	Link           string    `gorm:"column:link;primaryKey;not null;size:32"`
	Type           Type      `gorm:"column:captcha_type;"`
	Captcha        Captcha   `gorm:"column:captcha;not null;size:32"`
	SubmittedTime  time.Time `gorm:"column:submitted_time;"`
	ExpirationTime time.Time `gorm:"column:expiration_time;index;"`
}

func (model *CaptchaM) TableName() string {
	return "helix_captcha"
}
