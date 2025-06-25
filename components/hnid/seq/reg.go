package seq

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func InitSequenceIfNeed(code string, min uint64, max uint64, step uint64) error {
	exist, err := hdb.Exist[SequenceM](zplt.HelixPgDB().PG(), "code = ?", code)
	if err != nil {
		hlog.Err("hnid.seq.InitSequence:Exist", zap.String("code", code), zap.Error(err))
		return err
	}
	if exist {
		return nil
	}
	seqM := &SequenceM{
		Code:       code,
		MinStart:   min,
		MaxEnd:     max,
		Step:       step,
		CurrentSeq: 0,
		Version:    0,
	}
	err = hdb.Create[SequenceM](zplt.HelixPgDB().PG(), seqM)
	if err != nil {
		hlog.Err("hnid.seq.InitSequence:Create", zap.String("code", code), zap.Error(err))
		return err
	}
	return nil
}
