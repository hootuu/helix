package hvault

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/crypto/haes"
	"github.com/hootuu/hyle/crypto/hed25519"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	bufferSize     = 100
	prefixLen      = 20
	reloadInterval = 10 * 24 * time.Hour
	syncInterval   = 2 * time.Hour
)

type key struct {
	idx   []byte
	key   []byte
	usage uint64
}

var (
	gBuf            [bufferSize]*key
	gFill           = 0
	gLstSelectedIdx = 0
	gLstLoadTime    = time.UnixMilli(0)
	gLstSyncTime    = time.UnixMilli(0)
	gLock           sync.RWMutex
)

func doEncryptWithPwd(src []byte, pwdBytes []byte) ([]byte, error) {
	pwd := haes.Password(pwdBytes)
	coverSrc, err := pwd.Cover(src)
	if err != nil {
		hlog.Err("helix.vault.EncryptWithPwd", zap.Error(err))
		return nil, err
	}
	return Encrypt(coverSrc)
}

func doDecryptWithPwd(src []byte, pwdBytes []byte) ([]byte, error) {
	decSrc, err := Decrypt(src)
	if err != nil {
		return nil, err
	}
	pwd := haes.Password(pwdBytes)
	uncoverSrc, err := pwd.Uncover(decSrc)
	if err != nil {
		hlog.Err("helix.vault.DecryptWithPwd", zap.Error(err))
		return nil, err
	}
	return uncoverSrc, nil
}

func doEncrypt(src []byte) ([]byte, error) {
	selKey, err := doSelect()
	if err != nil {
		return nil, err
	}
	encBytes, err := haes.Encrypt(src, selKey.key)
	if err != nil {
		return nil, err
	}
	fullBytes := append(selKey.idx, encBytes...)
	return fullBytes, nil
}

func doDecrypt(src []byte) ([]byte, error) {
	if len(src) < prefixLen {
		hlog.Err("helix.vault.Decrypt: len(src) < prefixLen")
		return nil, errors.New("src len to short")
	}
	idxBytes := src[:prefixLen]
	encBytes := src[prefixLen:]

	var willUsePriKey []byte
	for _, wrap := range gBuf {
		if wrap == nil {
			continue
		}
		if bytes.Equal(wrap.idx, idxBytes) {
			willUsePriKey = wrap.key
		}
	}
	if willUsePriKey == nil {
		kcIdx := string(idxBytes)
		vaultM, err := hdb.Get[VaultM](zplt.HelixPgDB().PG(), "idx = ?", kcIdx)
		if err != nil {
			hlog.Err("helix.vault.Decrypt", zap.Error(err))
			return nil, err
		}
		willUsePriKey = vaultM.PrivateKey[:32]
	}
	decBytes, err := haes.Decrypt(encBytes, willUsePriKey)
	if err != nil {
		return nil, err
	}
	return decBytes, nil
}

func genNew() (*VaultM, error) {
	_, privateKey, err := hed25519.Random()
	if err != nil {
		hlog.Err("helix.hvault.genNew: hed25519.Random()", zap.Error(err))
		return nil, err
	}
	newKc := &VaultM{
		Idx:        idx.New(),
		PrivateKey: privateKey,
		Usage:      0,
		Available:  true,
	}
	return newKc, nil
}

func doSelect() (*key, error) {
	gLock.Lock()
	defer gLock.Unlock()
	if gFill == 0 || time.Now().Sub(gLstLoadTime) > reloadInterval {
		err := reload()
		if err != nil {
			return nil, err
		}
	}
	if gLstSelectedIdx == len(gBuf)-1 {
		gLstSelectedIdx = 0
	} else {
		gLstSelectedIdx++
	}
	curKw := gBuf[gLstSelectedIdx]
	curKw.usage++
	trySync()
	return curKw, nil
}

func trySync() {
	if time.Now().Sub(gLstSyncTime) < syncInterval {
		return
	}
	defer hlog.Elapse("helix.hvault.trySync", func() []zap.Field {
		return []zap.Field{zap.Time("gLstSyncTime", gLstSyncTime)}
	})()
	for _, item := range gBuf {
		err := updateUsage(string(item.idx), item.usage)
		if err != nil {
			hlog.Err("[ignore]helix.hvault.trySync: updateUsage", zap.Error(err))
			continue
		}
	}
	gLstSyncTime = time.Now()
}

func doLoad() error {
	defer hlog.Elapse("helix.hvault.doLoad", func() []zap.Field {
		return []zap.Field{zap.Time("gLstLoadTime", gLstLoadTime)}
	})()
	gFill = 0
	arr, err := loadAvailable(bufferSize)
	if err != nil {
		return err
	}
	for i, item := range arr {
		gBuf[i] = &key{
			[]byte(item.Idx),
			item.PrivateKey[:32],
			item.Usage,
		}
		gFill++
	}
	gLstLoadTime = time.Now()
	return nil
}

func reload() error {
	err := doLoad()
	if err != nil {
		return err
	}
	if gFill < bufferSize {
		err = batchGen(bufferSize - gFill)
		if err != nil {
			return err
		}
		err = doLoad()
		if err != nil {
			return err
		}
	}
	return nil
}

func batchGen(size int) error {
	defer hlog.Elapse("helix.hvault.batchGen", func() []zap.Field {
		return []zap.Field{zap.Int("size", size)}
	})()
	var arr []*VaultM
	for i := 0; i < size; i++ {
		newKc, err := genNew()
		if err != nil {
			return err
		}
		arr = append(arr, newKc)
	}
	err := multiCreate(arr)
	if err != nil {
		hlog.Err("helix.hvault.batchGen: multiCreate err", zap.Error(err))
		return err
	}
	return nil
}

func multiCreate(arr []*VaultM) error {
	err := hdb.MultiCreate[VaultM](zplt.HelixPgDB().PG(), arr)
	if err != nil {
		hlog.Err("helix.hvault.multiCreate", zap.Error(err))
		return nil
	}
	return nil
}

func updateUsage(idx string, usage uint64) error {
	mut := make(map[string]interface{})
	mut["usage"] = usage
	err := hdb.Update[VaultM](zplt.HelixPgDB().PG(), mut, "idx = ?", idx)
	if err != nil {
		hlog.Err("helix.hvault.updateUsage", zap.Error(err))
		return err
	}
	return nil
}

func loadAvailable(limit int) ([]*VaultM, error) {
	var arr []*VaultM
	arr, err := hdb.Find[VaultM](func() *gorm.DB {
		return zplt.HelixPgDB().PG().Model(&VaultM{}).
			Where("available = ?", true).Limit(limit)
	})
	if err != nil {
		hlog.Err("helix.hvault.loadAvailable", zap.Error(err))
		return nil, errors.New("loadAvailable error: " + err.Error())
	}
	return arr, nil
}
