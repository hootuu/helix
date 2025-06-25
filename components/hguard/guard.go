package hguard

import (
	"encoding/hex"
	"errors"
	"github.com/hootuu/helix/components/hvault"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hlocal"
	"github.com/hootuu/hyle/crypto/hed25519"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var gGuardPubLocalCache = hlocal.NewCache[[]byte](2*time.Hour, 24*time.Hour)

// GuardCreate : return PubKey, PriKey, error
func GuardCreate(biz string, alias string, call func(bizID string, pub []byte, pri []byte) error) error {
	pub, pri, err := hed25519.Random()
	if err != nil {
		return err
	}
	encryptPri, err := hvault.Encrypt(pri)
	if err != nil {
		hlog.Err("hguard.GuardCreate: hvault.Encrypt", zap.Error(err))
		return err
	}
	guardM := &GuardM{
		ID:           idx.New(),
		Biz:          biz,
		Alias:        alias,
		PubKey:       pub,
		PriKey:       encryptPri,
		Usage:        0,
		LstUsageTime: time.Now(),
	}
	err = hdb.Create[GuardM](zplt.HelixPgDB().PG(), guardM)
	if err != nil {
		hlog.Err("hguard.GuardCreate: hdb.Create", zap.Error(err))
		return err
	}
	gGuardPubLocalCache.Set(guardM.ID, &guardM.PubKey)
	return call(guardM.ID, pub, pri)
}

func GuardVerify(id string, data []byte, signature string) error {
	var pubKey []byte
	ptrPubKey, err := gGuardPubLocalCache.GetSet(id, func() (*[]byte, error) {
		guardM, err := hdb.Get[GuardM](zplt.HelixPgDB().PG(), "id = ?", id)
		if err != nil {
			hlog.Err("hguard.GuardVerify: hdb.Get", zap.Error(err))
			return nil, err
		}
		if guardM == nil {
			return nil, errors.New("no such guard for id: " + id)
		}
		return &(guardM.PubKey), nil
	})
	if err != nil {
		return err
	}
	pubKey = *ptrPubKey
	if len(pubKey) == 0 {
		hlog.Err("hguard.GuardVerify: len(pubKey) == 0")
		return errors.New("invalid guard public key")
	}

	bytesSign, err := hex.DecodeString(signature)
	if err != nil {
		hlog.Err("hguard.GuardVerify: hex.DecodeString(signature)",
			zap.String("id", id),
			zap.Error(err),
		)
		return err
	}
	valid := hed25519.Verify(pubKey, data, bytesSign)
	if !valid {
		return errors.New("guard verify: valid")
	}
	go func() {
		err := hdb.Update[GuardM](
			zplt.HelixPgDB().PG(),
			map[string]interface{}{
				"usage":          gorm.Expr("usage + 1"),
				"lst_usage_time": gorm.Expr("CURRENT_TIMESTAMP"),
			},
			"id = ?",
			id,
		)
		if err != nil {
			hlog.Err("[ignore]hguard.GuardVerify: hdb.Update",
				zap.String("id", id),
				zap.Error(err),
			)
		}
	}()
	return nil
}
