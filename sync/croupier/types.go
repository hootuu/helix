package croupier

import (
	"github.com/hootuu/hyle/crypto/hmd5"
	"github.com/hootuu/hyle/hypes/collar"
)

type ID = string

func BuildID(uniLink collar.Collar) ID {
	return hmd5.MD5(uniLink.Link().Str())
}
