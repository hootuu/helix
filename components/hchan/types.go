package hchan

import "github.com/hootuu/helix/components/htree"

type ID = htree.ID

const Root ID = 0

type Channel struct {
	ID       ID         `json:"id"`
	Name     string     `json:"name"`
	Icon     string     `json:"icon"`
	Seq      int        `json:"seq"`
	Children []*Channel `json:"children"`
}

func (c *Channel) AddChild(child *Channel) *Channel {
	c.Children = append(c.Children, child)
	return c
}
