package hnid

import (
	"github.com/hootuu/helix/components/hnid/seq"
	"github.com/hootuu/helix/components/hnid/slice"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func doNewLocalGenerator(code string, opt *Options) (Generator, error) {
	err := seq.InitSequenceIfNeed(code, opt.autoIncStart, opt.autoIncEnd, opt.autoIncStep)
	if err != nil {
		hlog.Err("hnid.doNewGenerator", zap.Error(err))
		return nil, err
	}
	autoSequence := seq.NewSequence(code)
	g := &localGenerator{
		bizSlice:       slice.NewNumber(opt.bizLen, opt.bizRef),
		timestampSlice: nil,
		autoIncSlice:   slice.NewAutoInc(opt.autoIncLen, autoSequence),
	}
	if opt.useTimestamp {
		g.timestampSlice = slice.NewTimestamp(opt.timestampType, opt.timestampUseDateFormat)
	}
	return g, nil
}

type localGenerator struct {
	timestampSlice *slice.Timestamp
	bizSlice       *slice.Number
	autoIncSlice   *slice.AutoInc
}

func (g *localGenerator) Next() NID {
	var id NID
	id.reset()
	if g.bizSlice != nil {
		id.Biz, id.BizLen = g.bizSlice.Build()
	}
	if g.timestampSlice != nil {
		id.Timestamp, id.TimestampLen = g.timestampSlice.Build()
	}
	if g.autoIncSlice != nil {
		id.AutoInc, id.AutoIncLen = g.autoIncSlice.Build()
	}
	return id
}

func (g *localGenerator) NextString() string {
	return g.Next().ToString()
}

func (g *localGenerator) NextUint64() uint64 {
	return g.Next().ToUint64()
}
