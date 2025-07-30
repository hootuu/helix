package hidem

import (
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hrds"
	"regexp"
	"sync"
	"time"
)

const (
	NoExpiration time.Duration = 0
)

type Factory interface {
	Check(idemCode string) (bool, error)
	MustCheck(idemCode string) error
}

var gFactoryMap = make(map[string]Factory)
var gFactoryMapMu sync.Mutex

func NewDbFactory(db *hdb.Database, code string, expiration time.Duration, cleanInterval time.Duration) (Factory, error) {
	if err := CheckFactoryCode(code); err != nil {
		return nil, err
	}
	gFactoryMapMu.Lock()
	defer gFactoryMapMu.Unlock()
	if _, ok := gFactoryMap[code]; ok {
		return nil, errors.New("repeated idempotent factory: " + code)
	}
	f, err := newDbFactory(db, code, expiration, cleanInterval)
	if err != nil {
		return nil, err
	}
	gFactoryMap[code] = f
	return f, nil
}

func NewCacheFactory(cache *hrds.Cache, code string, expiration time.Duration) (Factory, error) {
	if err := CheckFactoryCode(code); err != nil {
		return nil, err
	}
	gFactoryMapMu.Lock()
	defer gFactoryMapMu.Unlock()
	if _, ok := gFactoryMap[code]; ok {
		return nil, errors.New("repeated idempotent factory: " + code)
	}
	f, err := newCacheFactory(cache, code, expiration)
	if err != nil {
		return nil, err
	}
	gFactoryMap[code] = f
	return f, nil
}

func Each(call func(code string, g Factory)) {
	for code, f := range gFactoryMap {
		call(code, f)
	}
}

const gFactoryCodeRegexpTpl = `^[A-Za-z][A-Za-z0-9_]{0,32}$`

var gFactoryCodeRegexp = regexp.MustCompile(gFactoryCodeRegexpTpl)

func CheckFactoryCode(code string) error {
	matched := gFactoryCodeRegexp.MatchString(code)
	if !matched {
		return errors.New("invalid idem factory code[" + gFactoryCodeRegexpTpl + "]: " + code)
	}
	return nil
}

const gIdemCodeRegexpTpl = `^[A-Za-z0-9_.-:]{0,127}$`

var gIdemCodeRegexp = regexp.MustCompile(gIdemCodeRegexpTpl)

func CheckIdemCode(idemCode string) error {
	matched := gIdemCodeRegexp.MatchString(idemCode)
	if !matched {
		return errors.New("invalid idem code[" + gIdemCodeRegexpTpl + "]: " + idemCode)
	}
	return nil
}
