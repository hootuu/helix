package slice

import (
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

type Number struct {
	length         uint8
	referenceValue uint64
	useRand        bool
}

func NewNumber(l uint8, r uint64) *Number {
	if len(hcast.ToString(r)) > int(l) {
		hlog.Err("hnid.slice.NewNumber: len(hcast.ToString(r)) > int(l)",
			zap.Uint64("ref", r),
			zap.Uint8("len", l))
		r = 0
	}
	return &Number{
		length:         l,
		referenceValue: r,
	}
}

func (ns *Number) Build() (uint64, uint8) {
	if ns.length == 0 {
		return 0, 0
	} else if ns.length > 9 {
		hlog.Err("helix.hnid.slice.number.build: ns.length > 9")
		return 0, 0
	}
	return ns.referenceValue, ns.length
}
