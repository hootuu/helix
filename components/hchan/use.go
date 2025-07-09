package hchan

import (
	"github.com/hootuu/helix/helix"
)

func NewFactory(code string) (*Factory, error) {
	if err := helix.CheckCode(code); err != nil {
		return nil, err
	}
	cate, err := newFactory(code, 8, []uint{3, 3, 3, 3})
	if err != nil {
		return nil, err
	}
	helix.Use(cate.Helix())
	return cate, nil
}
