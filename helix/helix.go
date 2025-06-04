package helix

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyle/hync"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"regexp"
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

const gCodeRegexpTpl = `^[A-Za-z][A-Za-z0-9_]{0,32}$`

var gCodeRegexp = regexp.MustCompile(gCodeRegexpTpl)

func CheckCode(code string) error {
	matched := gCodeRegexp.MatchString(code)
	if !matched {
		return errors.New("invalid helix code[" + gCodeRegexpTpl + "]: " + code)
	}
	return nil
}

var gHelixBeenStartup = false
var gHelixArr []Helix
var gHelixMu sync.Mutex
var gHelixStartupSuccessOn = hync.NewOn()

func exist(code string) bool {
	for _, helix := range gHelixArr {
		if helix.code == code {
			return true
		}
	}
	return false
}

func doRegister(helix Helix) {
	gHelixMu.Lock()
	defer gHelixMu.Unlock()
	if bExist := exist(helix.code); bExist {
		hlog.Err("helix.doRegister: code repetition", zap.String("code", helix.code))
		hsys.Exit(errors.New("helix code repetition: " + helix.code))
	}
	gHelixArr = append(gHelixArr, helix)
	if gHelixBeenStartup {
		hsys.Warn("# Runtime startup the helix: [", helix.code, "] ...... #")
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

		hcfg.Dump(func(key string, val any) {
			if strings.Index(key, helix.code) > -1 {
				hsys.Info(" ** [", key, "] ==> ", val)
			}
		})
		hsys.Success("# Runtime startup the helix [", helix.code, "] [OK] #")
	}
}

func doStartup() {
	hsys.Warn("# Startup all registered helix ...... #")
	gHelixBeenStartup = true
	for _, helix := range gHelixArr {
		hsys.Warn("  ## Startup the helix: [", helix.code, "] ...... #")
		ctx, err := helix.startup()
		if err != nil {
			hlog.Err("start helix failed", zap.String("code", helix.code), zap.Error(err))
			hsys.Error("  ** Start helix exception: [", helix.code, "] #")
			return
		}
		helix.ctx = ctx
		if helix.ctx == nil {
			helix.ctx = context.Background()
		}
		hsys.Success("  ** Startup the helix [", helix.code, "] [OK] #")
	}
	hsys.Success("# Startup all registered helix [OK] #")

	hsys.Warn("# Display all init used configure items ...... #")
	hcfg.Dump(func(key string, val any) {
		if strings.Index(strings.ToLower(key), "password") > -1 ||
			strings.Index(strings.ToLower(key), "pwd") > -1 {
			hsys.Info("  ## [", key, "] ==>> ", "**********")
		} else {
			hsys.Info("  ## [", key, "] ==>> ", val)
		}
	})
	hsys.Success("# Display all init used configure items [OK] #")

	gHelixStartupSuccessOn.On()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println()
	hsys.Warn("# Shutdown the system ...... #")

	var cancelFuncArr []context.CancelFunc
	defer func() {
		if len(cancelFuncArr) > 0 {
			for _, cancelFunc := range cancelFuncArr {
				cancelFunc()
			}
		}
	}()
	for _, helix := range gHelixArr {
		ctx := helix.ctx
		if ctx == nil {
			ctx = context.Background()
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		cancelFuncArr = append(cancelFuncArr, cancel)
		hsys.Info("  ## Shutdown helix: [", helix.code, "] ......")
		helix.shutdown(ctx)
		hsys.Success("  ## Shutdown helix ", helix.code, " [OK]")
	}
	hsys.Success("# Shutting down the system [OK] #")
}
