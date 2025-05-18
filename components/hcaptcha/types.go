package hcaptcha

import (
	"errors"
	"github.com/hootuu/hyle/data/hcast"
	"time"
)

type Captcha = string

const (
	NilCaptcha Captcha = ""
)

type Type uint8

const (
	DigitsCaptcha      Type = 1
	LetterCaptcha      Type = 2
	LetterUpperCaptcha Type = 3
	LetterLowerCaptcha Type = 4
	MixtureCaptcha     Type = 5
)

type Options struct {
	Type       Type
	Length     int
	Expiration time.Duration
}

func NewOpt() *Options {
	return NewOptions(DigitsCaptcha, 4, 15*time.Minute)
}

func NewOptions(t Type, len int, exp time.Duration) *Options {
	return &Options{
		Type:       t,
		Length:     len,
		Expiration: exp,
	}
}

func (opt *Options) verify() error {
	switch opt.Type {
	case DigitsCaptcha,
		LetterCaptcha,
		LetterUpperCaptcha,
		LetterLowerCaptcha,
		MixtureCaptcha:
	default:
		return errors.New("invalid captcha type: " + hcast.ToString(opt.Type))
	}
	if opt.Length <= 0 || opt.Length > 32 {
		return errors.New("invalid captcha length: " + hcast.ToString(opt.Length))
	}
	if opt.Expiration < 0 {
		return errors.New("invalid captcha expiration: " + hcast.ToString(opt.Expiration))
	}
	return nil
}
