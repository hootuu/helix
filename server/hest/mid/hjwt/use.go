package hjwt

import (
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/herr"
)

type AuthenticateFunc func(dict dict.Dict) (interface{}, *herr.Error)
type RefreshFunc func(dict dict.Dict) (interface{}, *herr.Error)

func NewJwtMid(code string) *JwtMid {
	return nil //todo
}

func (mid *JwtMid) SetAuthenticate(path string, call AuthenticateFunc) {

}

func (mid *JwtMid) SetRefresh(path string, call RefreshFunc) {

}
