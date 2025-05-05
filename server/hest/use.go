package hest

import "github.com/hootuu/helix/helix"

func NewHest(code string) *Hest {
	h := newHest(code)
	helix.Use(h.Helix())
	return h
}
