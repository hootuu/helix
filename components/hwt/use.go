package hwt

import (
	"errors"
	"regexp"
)

func NewHwt(code string) (*Hwt, error) {
	if err := CheckCode(code); err != nil {
		return nil, err
	}
	return newHwt(code)
}

const gCodeRegexpTpl = `^[A-Za-z][A-Za-z0-9_]{0,16}$`

var gCodeRegexp = regexp.MustCompile(gCodeRegexpTpl)

func CheckCode(code string) error {
	matched := gCodeRegexp.MatchString(code)
	if !matched {
		return errors.New("invalid hwt code[" + gCodeRegexpTpl + "]: " + code)
	}
	return nil
}
