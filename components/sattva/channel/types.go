package channel

import (
	"fmt"
	"github.com/hootuu/hyle/crypto/hmd5"
	"github.com/hootuu/hyle/data/dict"
)

type Type int

type ID = string

func IdOf(t Type, code string) ID {
	return hmd5.MD5(fmt.Sprintf("%d:%s", t, code))
}

const IdNil = ""

type Link = string

type Config = dict.Dict

type Channel struct {
	Type  Type      `json:"type"`
	Code  string    `json:"code"`
	Link  Link      `json:"link"`
	Paras dict.Dict `json:"paras"`
}

type Builder interface {
	Default() Handler
	Build(id ID, cfg Config) (Handler, error)
}

type Handler interface {
	Wrap(chn *Channel) (*Channel, error)
	Identify(chn *Channel) (bool, error)
}
