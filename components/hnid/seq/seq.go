package seq

import (
	"errors"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"go.uber.org/zap"
	"sync"
)

type Sequence struct {
	code     string
	localSeq uint64
	start    uint64
	end      uint64
	lock     sync.Mutex
}

func NewSequence(code string) *Sequence {
	return &Sequence{
		code:     code,
		localSeq: 0,
		start:    0,
		end:      0,
	}
}

func (seq *Sequence) Next() uint64 {
	return seq.nextSequence()
}

func (seq *Sequence) refreshLocal() error {
	defer hlog.Elapse("helix.hnid.seq.Sequence.refreshLocal",
		func() []zap.Field {
			return []zap.Field{
				zap.String("code", seq.code),
				zap.Uint64("seq", seq.localSeq),
				zap.Uint64("start", seq.start),
				zap.Uint64("end", seq.end),
			}
		},
		func() []zap.Field {
			return []zap.Field{
				zap.String("code", seq.code),
				zap.Uint64("seq", seq.localSeq),
				zap.Uint64("start", seq.start),
				zap.Uint64("end", seq.end),
			}
		})()

	seqM, err := hpg.Get[SequenceM](zplt.HelixPgDB().PG(), "code = ?", seq.code)
	if err != nil {
		hlog.Err("helix.hnid.seq.PgSequence.refreshLocal: hpg.Get[SequenceM]", zap.Error(err))
		return err
	}
	if seqM == nil {
		hlog.Err("helix.hnid.seq.PgSequence.refreshLocal: hpg.Get[SequenceM] seqM == nil")
		return errors.New("no such seq data: code = " + seq.code)
	}

	localStart := seqM.CurrentSeq + 1
	localEnd := seqM.CurrentSeq + seqM.Step
	if localEnd > seqM.MaxEnd {
		localEnd = seqM.MaxEnd
	}

	dbNewCurrent := seqM.CurrentSeq + seqM.Step
	if dbNewCurrent >= seqM.MaxEnd {
		dbNewCurrent = seqM.MinStart
	}
	err = hpg.Update[SequenceM](zplt.HelixPgDB().PG(),
		map[string]any{
			"current_seq": dbNewCurrent,
			"version":     seqM.Version + 1,
		},
		"code = ? AND version = ?",
		seqM.Code, seqM.Version,
	)
	if err != nil {
		hlog.Err("helix.hnid.seq.PgSequence.refreshLocal: hpg.Update[SequenceM]", zap.Error(err))
		return err
	}

	seq.localSeq = localStart
	seq.start = localStart
	seq.end = localEnd

	return nil
}

func (seq *Sequence) checkWhetherNeedRefresh() bool {
	return (seq.end == 0) || seq.localSeq >= seq.end
}

func (seq *Sequence) nextSequence() uint64 {
	seq.lock.Lock()
	defer seq.lock.Unlock()
	if seq.checkWhetherNeedRefresh() {
		hretry.Universal(func() error {
			err := seq.refreshLocal()
			if err != nil {
				hlog.Err("helix.hnid.seq.PgSequence.nextSequence", zap.Error(err))
				return err
			}
			return nil
		})
	}
	seq.localSeq = seq.localSeq + 1
	return seq.localSeq
}
