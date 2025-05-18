package hwt

import (
	"github.com/hootuu/helix/helix"
)

func NewHwt(code string) (*Hwt, error) {
	if err := helix.CheckCode(code); err != nil {
		return nil, err
	}
	return newHwt(code)
}
