package helix

import (
	"errors"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
	"sync"
)

var gLoaderMap = make(map[string]*sync.Once)
var gLoaderMu sync.RWMutex

func OnceLoad(code string, build func()) {
	gLoaderMu.Lock()
	defer gLoaderMu.Unlock()
	var once *sync.Once
	once, ok := gLoaderMap[code]
	if !ok {
		once = &sync.Once{}
		gLoaderMap[code] = once
	}
	once.Do(func() {
		build()
	})
}

var gMustInitOnceMap = make(map[string]*sync.Once)
var gMustInitOnceMapMu sync.RWMutex

func MustInit(code string, doInit func() error) {
	if _, ok := gMustInitOnceMap[code]; ok {
		return
	}

	gMustInitOnceMapMu.Lock()
	defer gMustInitOnceMapMu.Unlock()
	var once *sync.Once
	once, ok := gMustInitOnceMap[code]
	if !ok {
		once = &sync.Once{}
		gMustInitOnceMap[code] = once
	}
	once.Do(func() {
		err := doInit()
		if err != nil {
			hlog.Err("helix.MustInit failed", zap.String("code", code), zap.Error(err))
			hsys.Exit(errors.New(fmt.Sprintf("init %s faild", code)))
		}
	})
}
