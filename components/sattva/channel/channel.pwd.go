package channel

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/hootuu/helix/components/hvault"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/crypto/hmd5"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

const (
	PwdEncryptPasswordKey = "encrypt_password"
	PwdPasswordKey        = "password"
)

type PwdBuilder struct {
	handler Handler
}

func NewPwdBuilder() *PwdBuilder {
	return &PwdBuilder{
		handler: newPwdHandler(nil),
	}
}

func (p *PwdBuilder) Default() Handler {
	return p.handler
}

func (p *PwdBuilder) Build(_ ID, cfg Config) (Handler, error) {
	if len(cfg) == 0 {
		return p.handler, nil
	}
	pwd := cfg.Get(PwdEncryptPasswordKey).String()
	if len(pwd) == 0 {
		return p.handler, nil
	}
	pwdBytes := hmd5.MD5Bytes(pwd)
	return newPwdHandler(pwdBytes), nil
}

type PwdHandler struct {
	encryptPwdBytes []byte
}

func newPwdHandler(encryptPwdBytes []byte) *PwdHandler {
	return &PwdHandler{
		encryptPwdBytes: encryptPwdBytes,
	}
}

const gPwdModLinkRegexpTpl = `^[a-zA-Z0-9_@.-]{1,32}$`

var gPwdModLinkRegexp = regexp.MustCompile(gPwdModLinkRegexpTpl)

func (h *PwdHandler) Wrap(input *Channel) (*Channel, error) {
	input.Link = strings.TrimSpace(input.Link)
	if matched := gPwdModLinkRegexp.MatchString(input.Link); !matched {
		return nil, errors.New("invalid link[" + gPwdModLinkRegexpTpl + "]: " + input.Link)
	}
	inputPwdBytes, err := pwdGetPwdFromChannel(input)
	if err != nil {
		return nil, err
	}

	var inputEncryptPwd []byte
	if h.encryptPwdBytes == nil {
		inputEncryptPwd, err = hvault.Encrypt(inputPwdBytes)
	} else {
		inputEncryptPwd, err = hvault.EncryptWithPwd(inputPwdBytes, h.encryptPwdBytes)
	}
	if err != nil {
		hlog.Err("sattva.channel.PwdHandler.Wrap: hvault.Encrypt", zap.Error(err))
		return nil, err
	}
	inputEncryptPwdHex := hex.EncodeToString(inputEncryptPwd)
	return &Channel{
		Type:  input.Type,
		Code:  input.Code,
		Link:  input.Link,
		Paras: dict.New(make(map[string]interface{})).Set(PwdPasswordKey, inputEncryptPwdHex),
	}, nil
}

func (h *PwdHandler) Identify(input *Channel) (bool, error) {
	chnID := IdOf(input.Type, input.Code)

	inputPwdBytes, err := pwdGetPwdFromChannel(input)
	if err != nil {
		return false, err
	}

	idChnM, err := hpg.Get[IdChannelM](zplt.HelixPgDB().PG(),
		"channel = ? AND link = ?",
		chnID, input.Link)
	if err != nil {
		hlog.Err("sattva.channel.PwdHandler.Identify", zap.Error(err))
		return false, err
	}
	if idChnM == nil {
		return false, errors.New("no such id channel [channel=" + input.Code + ", link=" + input.Link + "]")
	}
	ptrDbParas, err := hjson.FromBytes[dict.Dict](idChnM.Paras)
	if err != nil {
		hlog.Err("sattva.channel.PwdHandler.Identity: bytes to paras", zap.Error(err))
		return false, err
	}
	dbParas := *ptrDbParas
	if len(dbParas) == 0 {
		hlog.Err("sattva.channel.PwdHandler.Identity: invalid id channel paras", zap.Error(err))
		return false, errors.New("invalid id channel paras")
	}
	dbPwdEncryptHexStr := dbParas.Get(PwdPasswordKey).String()
	if len(dbPwdEncryptHexStr) == 0 {
		return false, errors.New("require valid password in db.channel.paras")
	}
	dbPwdEncryptBytes, err := hex.DecodeString(dbPwdEncryptHexStr)
	if err != nil {
		hlog.Err("sattva.channel.PwdHandler.Identity: hex.DecodeString", zap.Error(err))
		return false, err
	}
	var dbPwdDecryptBytes []byte
	if h.encryptPwdBytes == nil {
		dbPwdDecryptBytes, err = hvault.Decrypt(dbPwdEncryptBytes)
	} else {
		dbPwdDecryptBytes, err = hvault.DecryptWithPwd(dbPwdEncryptBytes, h.encryptPwdBytes)
	}
	if err != nil {
		hlog.Err("sattva.channel.PwdHandler.Identity: hvault.Decrypt", zap.Error(err))
		return false, err
	}

	if bytes.Equal(dbPwdDecryptBytes, inputPwdBytes) {
		return true, nil
	}
	return false, nil
}

func pwdGetPwdFromChannel(input *Channel) ([]byte, error) {
	if len(input.Paras) == 0 {
		return nil, errors.New("require channel.paras")
	}
	if _, ok := input.Paras[PwdPasswordKey]; !ok {
		return nil, errors.New("require channel.paras.password")
	}
	inputPwdStr := input.Paras.Get(PwdPasswordKey).String()
	if len(inputPwdStr) == 0 {
		return nil, errors.New("require valid channel.paras.password")
	}
	return hmd5.MD5Bytes(inputPwdStr), nil
}
