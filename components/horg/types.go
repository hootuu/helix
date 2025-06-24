package horg

import (
	"fmt"
	"github.com/hootuu/helix/components/hseq"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/hyle/hypes/collar"
)

type ID = htree.ID
type AuthID = hseq.ID

const Root ID = 0

const (
	CollarCode      = "horg"
	Collar4AuthCode = "horg_auth"
)

func Collar(id ID) collar.Collar {
	return collar.Build(CollarCode, fmt.Sprintf("%d", id))
}

func CollarAuth(id ID) collar.Collar {
	return collar.Build(Collar4AuthCode, fmt.Sprintf("%d", id))
}

type Organization struct {
	ID       ID              `json:"id"`
	Alias    string          `json:"alias"`
	Name     string          `json:"name"`
	Children []*Organization `json:"children"`
}

func (c *Organization) AddChild(child *Organization) *Organization {
	c.Children = append(c.Children, child)
	return c
}
