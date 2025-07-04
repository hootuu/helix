package hlink

import (
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hmath"
	"go.uber.org/zap"
	"math/rand/v2"
)

type Code = string

func newCodeNumbStr(seed uint64) string {
	seedStr := fmt.Sprintf("%d", seed)
	seedLen := len(seedStr)
	first := rand.Uint32N(10)
	lst := rand.Uint32N(10)
	lenStr := fmt.Sprintf("9860%d", seedLen)

	return fmt.Sprintf("%d%s%s%d", first, lenStr, seedStr, lst)
}

func newCode(codeNumbStr string) Code {
	c, err := hmath.Base10ToBase35(codeNumbStr)
	if err != nil {
		hlog.Err("helix.link.newCode err", zap.Error(err))
	}
	return c
}
