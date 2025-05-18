package hcaptcha

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func New(link string, opt *Options) (Captcha, error) {
	return newCaptcha(link, opt)
}

func Verify(link string, capCode Captcha) (bool, error) {
	return verifyCaptcha(link, capCode)
}

func init() {
	helix.Use(helix.BuildHelix("helix_captcha", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(&CaptchaM{})
		if err != nil {
			hlog.Err("hcaptcha.init: AutoMigrate", zap.Error(err))
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
