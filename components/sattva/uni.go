package sattva

import (
	"bytes"
	"errors"
	"github.com/hootuu/helix/components/sattva/channel"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"regexp"
)

const gChnCodeRegexpTpl = `^[A-Za-z][A-Za-z0-9_.]{0,32}$`

var gChnCodeRegexp = regexp.MustCompile(gChnCodeRegexpTpl)

func uniChannelRegister(chnType channel.Type, chnCode string, chnCfg channel.Config) (channel.ID, error) {
	matched := gChnCodeRegexp.MatchString(chnCode)
	if !matched {
		return channel.IdNil, errors.New("invalid channel code[" + gChnCodeRegexpTpl + "]: " + chnCode)
	}
	chnID := channel.IdOf(chnType, chnCode)
	chnM, err := hdb.Get[ChannelM](zplt.HelixPgDB().PG(), "id = ?", chnID)
	if err != nil {
		hlog.Err("sattva.uniChannelRegister: hdb.Get", zap.Error(err))
		return channel.IdNil, err
	}
	if chnM == nil {
		chnM = &ChannelM{
			ID:        chnID,
			Type:      chnType,
			Code:      chnCode,
			Config:    hjson.MustToBytes(chnCfg),
			Available: true,
		}
		err = hdb.Create[ChannelM](zplt.HelixPgDB().PG(), chnM)
		if err != nil {
			hlog.Err("sattva.uniChannelRegister: Create", zap.Error(err))
			return channel.IdNil, err
		}
		return chnM.ID, nil
	}

	newCfgBytes := hjson.MustToBytes(chnCfg)
	if chnM.Available == true && bytes.Equal(hjson.MustToBytes(chnM.Config), newCfgBytes) {
		return chnID, nil
	}

	channelUsed, err := hdb.Exist[channel.IdChannelM](zplt.HelixPgDB().PG(),
		"channel = ?", chnID)
	if err != nil {
		hlog.Err("sattva.uniChannelRegister: check used", zap.Error(err))
		return chnID, err
	}
	if channelUsed {
		return chnID, errors.New("there are already IDs that have used this channel. modification of the key information is not allowed")
	}

	mut := make(map[string]interface{})
	mut["config"] = newCfgBytes
	mut["available"] = true
	err = hdb.Update[ChannelM](zplt.HelixPgDB().PG(), mut, "id = ?", chnID)
	if err != nil {
		hlog.Err("sattva.uniChannelRegister: Update", zap.Error(err))
		return channel.IdNil, err
	}
	return chnM.ID, nil
}

func uniIdentificationCreate(chn *channel.Channel) (Identification, error) {
	chnID := channel.IdOf(chn.Type, chn.Code)
	handler, err := MustGetHandler(chn.Type, chnID)
	if err != nil {
		hlog.Err("sattva.doIdentificationCreate: MustGetHandler", zap.Error(err))
		return IdNil, err
	}
	wrapChn, err := handler.Wrap(chn)
	if err != nil {
		hlog.Err("sattva.doIdentificationCreate: handler.Wrap", zap.Error(err))
		return IdNil, err
	}
	idChannelExist, err := hdb.Exist[channel.IdChannelM](zplt.HelixPgDB().PG(),
		"channel = ? AND link = ?",
		chnID, chn.Link,
	)
	if err != nil {
		hlog.Err("sattva.doIdentificationCreate: exists id channel", zap.Any("channel", chn), zap.Error(err))
		return IdNil, err
	}
	if idChannelExist {
		return IdNil, errors.New("[" + chn.Code + "]" + chn.Link + " has been exists")
	}
	idM := &IdentificationM{
		ID: gUniIdGenerator.NextString(),
	}
	idChannelM := &channel.IdChannelM{
		ID:      idM.ID,
		Channel: chnID,
		Link:    wrapChn.Link,
		Paras:   hjson.MustToBytes(wrapChn.Paras),
	}
	err = hdb.Tx(zplt.HelixPgDB().PG(), func(tx *gorm.DB) error {
		if err := hdb.Create[IdentificationM](tx, idM); err != nil {
			return err
		}
		if err := hdb.Create[channel.IdChannelM](tx, idChannelM); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		hlog.Err("sattva.doIdentificationCreate: hdb.Tx", zap.Any("channel", chn), zap.Error(err))
		return IdNil, err
	}
	return idM.ID, nil
}

func uniIdentify(input *channel.Channel) (bool, error) {
	chnID := channel.IdOf(input.Type, input.Code)
	handler, err := MustGetHandler(input.Type, chnID)
	if err != nil {
		return false, err
	}
	return handler.Identify(input)
}
