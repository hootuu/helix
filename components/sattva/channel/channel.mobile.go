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
	PhoneKey        = "phone"
	gPhoneRegexpTpl = `^1[3-9][0-9]{9}$`
)

type MobileBuilder struct {
	handler Handler
}

func NewMobileBuilder() *MobileBuilder {
	return &MobileBuilder{
		handler: newMobileHandler(),
	}
}

func (p *MobileBuilder) Default() Handler {
	return p.handler
}

func (p *MobileBuilder) Build(_ ID, cfg Config) (Handler, error) {
	if len(cfg) == 0 {
		return p.handler, nil
	}
	return newMobileHandler(), nil
}

type MobileHandler struct {
}

func newMobileHandler() *MobileHandler {
	return &MobileHandler{}
}

func (h *MobileHandler) Wrap(input *Channel) (*Channel, error) {

	if len(input.Paras) == 0 {
		return nil, errors.New("channel parameters cannot be empty")
	}

	phone := input.Paras.Get(PhoneKey).String()
	if len(phone) == 0 {
		return nil, errors.New("phone parameter cannot be empty")
	}

	if matched := regexp.MustCompile(gPhoneRegexpTpl).MatchString(phone); !matched {
		return nil, errors.New("invalid phone format, must match " + gPhoneRegexpTpl)
	}

	return &Channel{
		Type:  input.Type,
		Code:  input.Code,
		Link:  input.Link,
		Paras: dict.New(make(map[string]interface{})).Set(PhoneKey, input.Paras.Get(PhoneKey).String()),
	}, nil
}

func (h *MobileHandler) Identify(input *Channel) (bool, string, error) {
	chnID := IdOf(input.Type, input.Code)

	idChnM, err := hdb.Get[IdChannelM](zplt.HelixPgDB().PG(),
		"channel = ? AND link = ?",
		chnID, input.Link)
	if err != nil {
		hlog.Err("sattva.channel.MobileHandler.Identify", zap.Error(err))
		return false, IdNil, err
	}
	if idChnM == nil {
		return false, IdNil, errors.New("no such id channel [channel=" + input.Code + ", link=" + input.Link + "]")
	}

	ptrDbParas, err := hjson.FromBytes[dict.Dict](idChnM.Paras)
	if err != nil {
		hlog.Err("sattva.channel.MobileHandler.Identity: bytes to paras", zap.Error(err))
		return false, IdNil, err
	}

	dbParas := *ptrDbParas
	if len(dbParas) == 0 {
		hlog.Err("sattva.channel.MobileHandler.Identity: invalid id channel paras", zap.Error(err))
		return false, IdNil, errors.New("invalid id channel paras")
	}

	dbPhone := dbParas.Get(PhoneKey).String()
	if len(dbPhone) == 0 {
		hlog.Err("sattva.channel.MobileHandler.Identity: invalid id channel phone", zap.Error(err))
		return false, IdNil, errors.New("invalid id channel phone")
	}

	if dbPhone != input.Paras.Get(PhoneKey).String() {
		hlog.Err("sattva.channel.MobileHandler.Identity: phone mismatch",
			zap.String("dbPhone", dbPhone),
			zap.String("inputPhone", input.Paras.Get(PhoneKey).String()))
		return false, IdNil, errors.New("phone mismatch")
	}

	return true, idChnM.ID, nil
}
