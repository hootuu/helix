package helix

import (
	"context"
	"errors"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type Helix struct {
	code     string
	startup  func() (context.Context, error)
	shutdown func(ctx context.Context)
	ctx      context.Context
}

func BuildHelix(
	code string,
	startup func() (context.Context, error),
	shutdown func(ctx context.Context),
) Helix {
	return Helix{
		code:     code,
		startup:  startup,
		shutdown: shutdown,
		ctx:      nil,
	}
}

var gHelixBeenStartup = false
var gHelixMap = make(map[string]Helix)
var gHelixMu sync.Mutex

func doRegister(helix Helix) {
	gHelixMu.Lock()
	defer gHelixMu.Unlock()
	if _, ok := gHelixMap[helix.code]; ok {
		hlog.Err("helix.doRegister: code repetition")
		hsys.Exit(errors.New("helix code repetition"))
	}
	gHelixMap[helix.code] = helix
	if gHelixBeenStartup {
		hsys.Info("\n# Runtime startup the helix: [", helix.code, "] ...... #")
		ctx, err := helix.startup()
		if err != nil {
			hlog.Err("runtime start helix failed", zap.String("code", helix.code), zap.Error(err))
			hsys.Error("# Runtime Start helix exception: [", helix.code, "] #")
			return
		}
		helix.ctx = ctx
		if helix.ctx == nil {
			helix.ctx = context.Background()
		}
		hsys.Success("# Runtime startup the helix [", helix.code, "] [OK] #\n")
	}
}

func doStartup() {
	hsys.Info("\n# Startup all registered helix ...... #")
	for code, helix := range gHelixMap {
		hsys.Info("  ## Startup the helix: [", code, "] ...... #")
		ctx, err := helix.startup()
		if err != nil {
			hlog.Err("start helix failed", zap.String("code", code), zap.Error(err))
			hsys.Error("  # Start helix exception: [", code, "] #")
			return
		}
		helix.ctx = ctx
		if helix.ctx == nil {
			helix.ctx = context.Background()
		}
		hsys.Success("  ## Startup the helix [", code, "] [OK] #")
	}
	hsys.Success("# Startup all registered helix [OK] #\n")

	hsys.Info("\n# Display all used configure items ...... #")
	hcfg.Dump(func(key string, val any) {
		if strings.Index(strings.ToLower(key), "password") > -1 ||
			strings.Index(strings.ToLower(key), "pwd") > -1 {
			hsys.Info("  ## [", key, "] ==>> ", "**********")
		} else {
			hsys.Info("  ## [", key, "] ==>> ", val)
		}
	})
	hsys.Success("# Display all used configure items [OK] #\n")

	gHelixBeenStartup = true

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	hsys.Warn("\n# Shutdown the system ...... #")

	var cancelFuncArr []context.CancelFunc
	defer func() {
		if len(cancelFuncArr) > 0 {
			for _, cancelFunc := range cancelFuncArr {
				cancelFunc()
			}
		}
	}()
	for code, helix := range gHelixMap {
		ctx, cancel := context.WithTimeout(helix.ctx, 5*time.Second)
		cancelFuncArr = append(cancelFuncArr, cancel)
		hsys.Info("  ## Shutdown helix: [", code, "] ......")
		helix.shutdown(ctx)
		hsys.Success("  ## Shutdown helix ", code, " [OK]")
	}
	hsys.Success("# Shutting down the system [OK] #")
}
