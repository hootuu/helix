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
	DeviceIDKey        = "device_id"
	gDeviceIDRegexpTpl = `^[A-Za-z0-9]{6,64}$`
)

type DeviceBuilder struct {
	handler Handler
}

func NewDeviceBuilder() *DeviceBuilder {
	return &DeviceBuilder{
		handler: newDeviceHandler(),
	}
}

func (p *DeviceBuilder) Default() Handler {
	return p.handler
}

func (p *DeviceBuilder) Build(_ ID, cfg Config) (Handler, error) {
	if len(cfg) == 0 {
		return p.handler, nil
	}
	return newDeviceHandler(), nil
}

type DeviceHandler struct {
}

func newDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

func (h *DeviceHandler) Wrap(input *Channel) (*Channel, error) {

	if len(input.Paras) == 0 {
		return nil, errors.New("channel parameters cannot be empty")
	}

	deviceID := input.Paras.Get(DeviceIDKey).String()
	if len(deviceID) == 0 {
		return nil, errors.New("device_id parameter cannot be empty")
	}

	if matched := regexp.MustCompile(gDeviceIDRegexpTpl).MatchString(deviceID); !matched {
		return nil, errors.New("invalid device_id[" + gDeviceIDRegexpTpl + "]: " + deviceID)
	}

	return &Channel{
		Type:  input.Type,
		Code:  input.Code,
		Link:  input.Link,
		Paras: dict.New(make(map[string]interface{})).Set(DeviceIDKey, input.Paras.Get(DeviceIDKey).String()),
	}, nil
}

func (h *DeviceHandler) Identify(input *Channel) (bool, string, error) {
	chnID := IdOf(input.Type, input.Code)

	idChnM, err := hdb.Get[IdChannelM](zplt.HelixPgDB().PG(),
		"channel = ? AND link = ?",
		chnID, input.Link)
	if err != nil {
		hlog.Err("sattva.channel.DeviceHandler.Identify", zap.Error(err))
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
		hlog.Err("sattva.channel.DeviceHandler.Identity: bytes to paras", zap.Error(err))
		return false, IdNil, err
	}
	dbParas := *ptrDbParas
	if len(dbParas) == 0 {
		hlog.Err("sattva.channel.DeviceHandler.Identity: invalid id channel paras", zap.Error(err))
		return false, IdNil, errors.New("invalid id channel paras")
	}

	dbDeviceID := dbParas.Get(DeviceIDKey).String()
	if len(dbDeviceID) == 0 {
		hlog.Err("sattva.channel.DeviceHandler.Identity: invalid id channel device id", zap.Error(err))
		return false, IdNil, errors.New("invalid id channel device id")
	}

	if dbDeviceID != input.Paras.Get(DeviceIDKey).String() {
		hlog.Err("sattva.channel.DeviceHandler.Identity: device id not match",
			zap.String("dbDeviceID", dbDeviceID),
			zap.String("inputDeviceID", input.Paras.Get(DeviceIDKey).String()),
		)
		return false, IdNil, errors.New("device id not match")
	}

	return true, idChnM.ID, nil
}
