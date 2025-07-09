package hlink

import (
	"context"
	"errors"
	"github.com/hootuu/helix/components/hseq"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/crypto/hmd5"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"go.uber.org/zap"
	"strings"
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

func Validate(biz string, code Code) (bool, error) {
	tx := zplt.HelixDB().DB()
	b, err := hdb.Exist[LinkCodeM](tx, "code = ? AND biz = ?", code, biz)
	if err != nil {
		return false, err
	}
	return b, nil
}

func GetMajor(biz string, code Code) (collar.Collar, error) {
	tx := zplt.HelixDB().DB()
	codeM, err := hdb.Get[LinkCodeM](tx, "code = ? AND biz = ?", code, biz)
	if err != nil {
		return "", err
	}
	if codeM == nil {
		return "", errors.New("code not found")
	}
	return collar.FromID(codeM.Link)
}

func Bind(
	ctx context.Context,
	biz string,
	code Code,
	relation string,
	counterpart collar.Collar,
) error {
	tx := zplt.HelixPgCtx(ctx)
	codeM, err := hdb.MustGet[LinkCodeM](tx, "code = ? AND biz = ?", code, biz)
	if err != nil {
		return err
	}
	counterpartStr := counterpart.ToID()
	linkID := hmd5.MD5(strings.Join([]string{
		biz,
		codeM.Link,
		relation,
		counterpartStr,
	}, "-"))
	b, err := hdb.Exist[LinkM](tx, "id = ?", linkID)
	if err != nil {
		return err
	}
	if b {
		return errors.New("the relation already exists")
	}
	linkM := &LinkM{
		ID:          linkID,
		Biz:         codeM.Biz,
		Major:       codeM.Link,
		Relation:    relation,
		Counterpart: counterpartStr,
	}
	err = hdb.Create[LinkM](tx, linkM)
	if err != nil {
		return err
	}
	return nil
}

func Unbind(
	ctx context.Context,
	biz string,
	major collar.Collar,
	relation string,
	counterpart collar.Collar,
) error {
	tx := zplt.HelixPgCtx(ctx)
	majorStr := major.ToID()
	counterpartStr := counterpart.ToID()
	linkID := hmd5.MD5(strings.Join([]string{
		biz,
		majorStr,
		relation,
		counterpartStr,
	}, "-"))
	err := hdb.Delete[LinkM](tx, "id = ?", linkID)
	if err != nil {
		return err
	}
	return nil
}
