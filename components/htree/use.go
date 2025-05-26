package htree

import "github.com/hootuu/helix/helix"

func NewTree(code string, version uint, cfg []uint) (*Tree, error) {
	if err := helix.CheckCode(code); err != nil {
		return nil, err
	}
	t, err := newTree(code, version, cfg)
	if err != nil {
		return nil, err
	}
	helix.Use(t.Helix())
	return t, nil
}
