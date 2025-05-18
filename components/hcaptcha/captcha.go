package hcaptcha

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/crypto/hmd5"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hync"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func newCaptcha(link string, opt *Options) (Captcha, error) {
	defer cleanExpiration()
	linkID := hmd5.MD5(link)
	if err := opt.verify(); err != nil {
		return NilCaptcha, err
	}
	cptCode := Generate(opt)
	current := time.Now()
	capM := &CaptchaM{
		Link:           linkID,
		Type:           opt.Type,
		Captcha:        cptCode,
		SubmittedTime:  current,
		ExpirationTime: current.Add(opt.Expiration),
	}
	err := hpg.Create[CaptchaM](zplt.HelixPgDB().PG(), capM)
	if err != nil {
		hlog.Err("hcaptcha.newCaptcha: hpg.Create", zap.Error(err))
		return NilCaptcha, err
	}
	return cptCode, nil
}

func verifyCaptcha(link string, capCode Captcha) (bool, error) {
	linkID := hmd5.MD5(link)
	capM, err := hpg.Get[CaptchaM](zplt.HelixPgDB().PG(), "link = ?", linkID)
	if err != nil {
		hlog.Err("hcaptcha.verifyCaptcha: hpg.Get", zap.Error(err))
		return false, err
	}
	if capM == nil {
		return false, nil
	}
	if time.Now().After(capM.ExpirationTime) {
		return false, nil
	}
	if capM.Captcha != capCode {
		return false, nil
	}
	return true, nil
}

func doCleanExpiration() {
	beforeDays := hcfg.GetInt("helix.captcha.clean.before.days", 1)
	tx := zplt.HelixPgDB().PG().Unscoped().
		Where("expiration_time < CURRENT_TIMESTAMP - ?", gorm.Expr("? * INTERVAL '1 DAY'", beforeDays)).
		Delete(&CaptchaM{})
	if tx.Error != nil {
		hlog.Err("[ignore]hcaptcha.verifyCaptcha: Delete", zap.Error(tx.Error))
		return
	}
	if tx.RowsAffected > 0 {
		hlog.Logger().Info("hcaptcha.verifyCaptcha: Delete", zap.Int64("RowsAffected", tx.RowsAffected))
	}
}

var gCleanLine = hync.NewLine()
var gCleanLstTime = time.UnixMilli(0)

func cleanExpiration() {
	_ = gCleanLine.Do(func() error {
		interval := hcfg.GetDuration("helix.captcha.clean.interval", 6*time.Hour)
		current := time.Now()
		if current.Sub(gCleanLstTime) < interval {
			return nil
		}
		gCleanLstTime = current
		go doCleanExpiration()
		return nil
	})
}
