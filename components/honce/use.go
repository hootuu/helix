package honce

import (
	"errors"
	"regexp"
)

func Do(strCode string, call func() error) error {
	if err := CheckOnceCode(strCode); err != nil {
		return err
	}
	return doOnce(strCode, call)
}

var gOnceCodeRegexpTpl = `^[0-9][A-Za-z0-9_.-]{0,127}$`
var gOnceCodeRegexp = regexp.MustCompile(gOnceCodeRegexpTpl)

func CheckOnceCode(onceCode string) error {
	matched := gOnceCodeRegexp.MatchString(onceCode)
	if !matched {
		return errors.New("invalid once code[" + gOnceCodeRegexpTpl + "]: " + onceCode)
	}
	return nil
}
