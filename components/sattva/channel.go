package sattva

import (
	"errors"
	"github.com/hootuu/helix/components/sattva/channel"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

const (
	handlerAnyID = "*"
)

type localHandler struct {
	handler     channel.Handler
	lstSyncTime time.Time
}

var gFactoryBuilders = make(map[channel.Type]channel.Builder)
var gFactoryHandlers = make(map[channel.Type]map[channel.ID]*localHandler)
var gFactoryMu sync.Mutex
var gFactoryLstSyncTime = time.UnixMilli(0)
var gFactorySyncing = false

func RegisterBuilder(t channel.Type, builder channel.Builder) {
	gFactoryMu.Lock()
	defer gFactoryMu.Unlock()
	if _, ok := gFactoryBuilders[t]; ok {
		hlog.Err("sattva.channel.RegisterBuilder: builder repeated",
			zap.Int("type", int(t)))
		return
	}
	gFactoryBuilders[t] = builder
	handlerMap := make(map[channel.ID]*localHandler)
	defaultHandler := builder.Default()
	if defaultHandler != nil {
		handlerMap[handlerAnyID] = &localHandler{
			handler:     defaultHandler,
			lstSyncTime: time.Now(),
		}
	}
	gFactoryHandlers[t] = handlerMap
}

func MustGetHandler(t channel.Type, id channel.ID) (channel.Handler, error) {
	localHandlersReload()
	handlerMap, ok := gFactoryHandlers[t]
	if !ok {
		return nil, errors.New("no such handler for: [" + hcast.ToString(t) + "]")
	}
	handler, ok := handlerMap[id]
	if !ok {
		handler, ok = handlerMap[handlerAnyID]
		if !ok {
			return nil, errors.New("no such handler for: [" + hcast.ToString(t) + "]" + id)
		}
	}
	return handler.handler, nil
}

func localHandlersReload() {
	syncInterval := hcfg.GetDuration("sattva.channel.sync.interval", 30*time.Minute)
	if time.Now().Sub(gFactoryLstSyncTime) < syncInterval {
		return
	}
	gFactoryMu.Lock()
	defer gFactoryMu.Unlock()
	if !gFactorySyncing {
		gFactorySyncing = true
		func() {
			syncSuccess := false
			defer func() {
				gFactorySyncing = false
				if syncSuccess {
					gFactoryLstSyncTime = time.Now()
				}
			}()
			chnArr, err := hpg.Find[ChannelM](func() *gorm.DB {
				return zplt.HelixPgDB().PG().
					Where("available = ?", true)
			})
			if err != nil {
				hlog.Err("[ignore]sattva.localHandlersReload: find channel", zap.Error(err))
				return
			}
			for _, chnM := range chnArr {
				builder := doGetBuilder(chnM.Type)
				if builder == nil {
					hlog.Err("[ignore]sattva.localHandlersReload: no such builder",
						zap.Int("type", int(chnM.Type)))
					continue
				}
				handlerMap, ok := gFactoryHandlers[chnM.Type]
				if !ok {
					handlerMap = make(map[channel.ID]*localHandler)
					gFactoryHandlers[chnM.Type] = handlerMap
				}
				localH, ok := handlerMap[chnM.ID]
				needRebuild := !ok || localH.lstSyncTime.After(chnM.UpdatedAt)
				if needRebuild {
					ptrCfg, err := hjson.FromBytes[channel.Config](chnM.Config)
					if err != nil {
						hlog.Err("[ignore]sattva.localHandlersReload: parse channel config", zap.Error(err))
						continue
					}
					handler, err := builder.Build(chnM.ID, *ptrCfg)
					if err != nil {
						hlog.Err("[ignore]sattva.localHandlersReload: builder.Builder", zap.Error(err))
						continue
					}
					localH = &localHandler{
						handler:     handler,
						lstSyncTime: chnM.UpdatedAt,
					}
					handlerMap[chnM.ID] = localH
				}
			}
			syncSuccess = true
		}()
	}
}

func doGetBuilder(t channel.Type) channel.Builder {
	builder, ok := gFactoryBuilders[t]
	if !ok {
		return nil
	}
	return builder
}
