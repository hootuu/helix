package hlink

import (
	"context"
	"github.com/hootuu/helix/components/hseq"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"go.uber.org/zap"
)

func Generate(ctx context.Context, biz string, major collar.Collar) (c Code, err error) {
	defer hlog.Elapse("helix.link.Generate", func() []zap.Field {
		return nil
	}, func() []zap.Field {
		if err != nil {
			return []zap.Field{
				zap.String("biz", biz),
				zap.String("major", major.ToString()),
				zap.Error(err),
			}
		}
		return nil
	})()
	seed, err := hseq.Next(ctx, collar.Build("helix_link", biz))
	if err != nil {
		return "", err
	}
	c = newCode(newCodeNumbStr(uint64(seed)))
	tx := zplt.HelixPgCtx(ctx)
	codeM := &LinkCodeM{
		Link: major.ToID(),
		Biz:  biz,
		Code: c,
	}
	err = hdb.Create[LinkCodeM](tx, codeM)
	if err != nil {
		return "", nil
	}
	return c, nil
}

func Validate(code Code) error {
	return nil
}

func GetMajor(code Code) (collar.Collar, error) {
	return "", nil
}

func Bind(code Code, relation string, counterpart collar.Collar) error {
	return nil
}

func Unbind(code Code, relation string, counterpart collar.Collar) error {
	return nil
}
