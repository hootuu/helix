package sattva

import "github.com/hootuu/hyle/data/dict"

type AuthType string

type AuthInfo struct {
	AuthType AuthType
	Paras    dict.Dict
}

type Authenticator interface {
	GetType() AuthType
	Authenticate(dict dict.Dict) (bool, error)
}
