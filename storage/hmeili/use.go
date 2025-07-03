package hmeili

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hlog"
	"sync"
)

var gMap = make(map[string]*Meili)
var gMutex sync.Mutex

func Register(code string) {
	gMutex.Lock()
	defer gMutex.Unlock()
	if _, ok := gMap[code]; ok {
		hlog.Err("hmeili.doRegister: meili repetition")
		return
	}
	m := New(code)
	gMap[code] = m
	helix.Use(m.Helix())
}

func GetMeili(code string) *Meili {
	gMutex.Lock()
	defer gMutex.Unlock()
	m, ok := gMap[code]
	if !ok {
		hlog.Err("hmeili.GetMeili: meili nil")
		return nil
	}
	return m
}
