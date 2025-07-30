package helix

import (
	"context"
	"errors"
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

var gLogger = hlog.GetLogger("helix")

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
		gLogger.Error("Helix already registered, will exit", zap.String("code", helix.code))
		hsys.Exit(errors.New("helix code repetition: " + helix.code))
	}
	gHelixArr = append(gHelixArr, helix)
	if gHelixBeenStartup {
		gLogger.Info("helix:["+helix.code+"] starting...", zap.String("code", helix.code))
		hsys.Warn("# Runtime startup the helix: [", helix.code, "] ...... #")
		ctx, err := helix.startup()
		if err != nil {
			gLogger.Error("runtime start helix failed, will exit", zap.String("code", helix.code), zap.Error(err))
			hsys.Error("# Runtime Start helix exception: [", helix.code, "] #")
			hsys.Exit(err)
			return
		}
		helix.ctx = ctx
		if helix.ctx == nil {
			helix.ctx = context.Background()
		}

		hcfg.Dump(func(key string, val any) {
			if strings.Index(key, helix.code) > -1 {
				gLogger.Info("helix:["+helix.code+"] cfg", zap.String("key", key), zap.Any("val", val))
				hsys.Info(" ** [", key, "] ==> ", val)
			}
		})
		gLogger.Info("helix:["+helix.code+"] cfg", zap.String("code", helix.code))
		hsys.Success("# Runtime startup the helix [", helix.code, "] [OK] #")
	}
}

func doStartup() {
	gLogger.Info("startup all registered helix ......")
	hsys.Warn("# Startup all registered helix ...... #")
	gHelixBeenStartup = true
	for _, helix := range gHelixArr {
		gLogger.Info("helix:["+helix.code+"] starting...", zap.String("code", helix.code))
		hsys.Warn("  ## Startup the helix: [", helix.code, "] ...... #")
		ctx, err := helix.startup()
		if err != nil {
			gLogger.Error("runtime start helix failed", zap.String("code", helix.code), zap.Error(err))
			hsys.Error("  ** Start helix exception: [", helix.code, "] #")
			return
		}
		helix.ctx = ctx
		if helix.ctx == nil {
			helix.ctx = context.Background()
		}
		gLogger.Info("helix:["+helix.code+"] start OK", zap.String("code", helix.code))
		hsys.Success("  ** Startup the helix [", helix.code, "] [OK] #")
	}
	gLogger.Info("startup all registered helix OK")
	hsys.Success("# Startup all registered helix [OK] #")

	gLogger.Info("# Display all init used configure items ...... #")
	hsys.Warn("# Display all init used configure items ...... #")
	hcfg.Dump(func(key string, val any) {
		if strings.Index(strings.ToLower(key), "password") > -1 ||
			strings.Index(strings.ToLower(key), "pwd") > -1 {
			gLogger.Info("  ## [" + key + "] ==>> **********")
			hsys.Info("  ## [", key, "] ==>> ", "**********")
		} else {
			gLogger.Info("  ## ["+key+"] ==>> ", zap.String("key", key), zap.Any("val", val))
			hsys.Info("  ## [", key, "] ==>> ", val)
		}
	})
	gLogger.Info("# Display all init used configure items [OK] #")
	hsys.Success("# Display all init used configure items [OK] #")

	gHelixStartupSuccessOn.On()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gLogger.Info("# Shutdown the system ...... #")
	hsys.Warn("\n\n# Shutdown the system ...... #")

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
		gLogger.Info("## Shutdown helix: [" + helix.code + "] ......")
		hsys.Info("  ## Shutdown helix: [", helix.code, "] ......")
		helix.shutdown(ctx)
		gLogger.Info("## Shutdown helix " + helix.code + " [OK]")
		hsys.Success("  ## Shutdown helix ", helix.code, " [OK]")
	}
	gLogger.Info("#  Shutting down the system [OK] [OK]")
	hsys.Success("# Shutting down the system [OK] #")
}
