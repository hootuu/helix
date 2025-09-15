package channel

import (
	"errors"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"regexp"
)

const (
	OpenIDKey        = "openid"
	gOpenIDRegexpTpl = `^[A-Za-z0-9]{6,64}$`
)

type WechatBuilder struct {
	handler Handler
}

func NewWechatBuilder() *WechatBuilder {
	return &WechatBuilder{
		handler: newWechatHandler(),
	}
}

func (p *WechatBuilder) Default() Handler {
	return p.handler
}

func (p *WechatBuilder) Build(_ ID, cfg Config) (Handler, error) {
	if len(cfg) == 0 {
		return p.handler, nil
	}
	return newWechatHandler(), nil
}

type WechatHandler struct {
}

func newWechatHandler() *WechatHandler {
	return &WechatHandler{}
}

func (h *WechatHandler) Wrap(input *Channel) (*Channel, error) {

	if len(input.Paras) == 0 {
		return nil, errors.New("channel parameters cannot be empty")
	}

	openID := input.Paras.Get(OpenIDKey).String()
	if len(openID) == 0 {
		return nil, errors.New("open_id parameter cannot be empty")
	}

	if matched := regexp.MustCompile(gOpenIDRegexpTpl).MatchString(openID); !matched {
		return nil, errors.New("invalid open_id[" + gOpenIDRegexpTpl + "]: " + openID)
	}

	return &Channel{
		Type:  input.Type,
		Code:  input.Code,
		Link:  input.Link,
		Paras: dict.New(make(map[string]interface{})).Set(OpenIDKey, input.Paras.Get(OpenIDKey).String()),
	}, nil
}

func (h *WechatHandler) Identify(input *Channel) (bool, string, error) {
	chnID := IdOf(input.Type, input.Code)

	idChnM, err := hdb.Get[IdChannelM](zplt.HelixPgDB().PG(),
		"channel = ? AND link = ?",
		chnID, input.Link)
	if err != nil {
		hlog.Err("sattva.channel.WechatHandler.Identify", zap.Error(err))
		return false, IdNil, err
	}
	if idChnM == nil || idChnM.ID == "" {
		return false, IdNil, errors.New("no such id channel [channel=" + input.Code + ", link=" + input.Link + "]")
	}

	if idChnM.ID == "" {
		return false, IdNil, errors.New("id channel is nil [channel=" + input.Code + ", link=" + input.Link + "]")
	}

	ptrDbParas, err := hjson.FromBytes[dict.Dict](idChnM.Paras)
	if err != nil {
		hlog.Err("sattva.channel.WechatHandler.Identity: bytes to paras", zap.Error(err))
		return false, IdNil, err
	}
	dbParas := *ptrDbParas
	if len(dbParas) == 0 {
		hlog.Err("sattva.channel.WechatHandler.Identity: invalid id channel paras", zap.Error(err))
		return false, IdNil, errors.New("invalid id channel paras")
	}

	dbOpenID := dbParas.Get(OpenIDKey).String()
	if len(dbOpenID) == 0 {
		hlog.Err("sattva.channel.WechatHandler.Identity: invalid id channel open id", zap.Error(err))
		return false, IdNil, errors.New("invalid id channel open id")
	}

	if dbOpenID != input.Paras.Get(OpenIDKey).String() {
		hlog.Err("sattva.channel.WechatHandler.Identity: open id not match",
			zap.String("dbOpenID", dbOpenID),
			zap.String("inputOpenID", input.Paras.Get(OpenIDKey).String()),
		)
		return false, IdNil, errors.New("open id not match")
	}

	return true, idChnM.ID, nil
}
