package hest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/server/hest/mid"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hsys"
	"net/http"
	"strings"
	"time"
)

type Hest struct {
	code       string
	httpServer *http.Server
	ginEngine  *gin.Engine
}

func newHest(code string) *Hest {
	h := &Hest{code: code}
	h.init()
	return h
}

func (h *Hest) Router(call func(router *gin.RouterGroup)) {
	call(&h.ginEngine.RouterGroup)
}

func (h *Hest) Helix() helix.Helix {
	return helix.BuildHelix(h.code, h.startup, h.shutdown)
}

func (h *Hest) init() {
	ginEngine := gin.New()
	ginEngine.Use(mid.NewLoggerMid(h.code).Handle)
	ginEngine.Use(mid.ErrHandle)
	cfgCode := strings.ToLower(h.code)
	addr := hcfg.GetString(fmt.Sprintf("hest.%s.addr", cfgCode), ":9860")
	readTimeout := hcfg.GetDuration(fmt.Sprintf("rest.%s.read.timeout", cfgCode), 30*time.Second)
	readHeaderTimeout := hcfg.GetDuration(fmt.Sprintf("rest.%s.read.header.timeout", cfgCode), 30*time.Second)
	writeTimeout := hcfg.GetDuration(fmt.Sprintf("rest.%s.write.timeout", cfgCode), 30*time.Second)
	idleTimeout := hcfg.GetDuration(fmt.Sprintf("rest.%s.idle.timeout", cfgCode), 30*time.Second)

	h.httpServer = &http.Server{
		Handler:           ginEngine,
		Addr:              addr,
		ReadTimeout:       readTimeout * time.Second,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
		IdleTimeout:       idleTimeout * time.Second,
	}

	h.ginEngine = ginEngine
}

func (h *Hest) startup() (context.Context, error) {
	hsys.Info("\n# Hest [", h.code, "] startup on: ", h.httpServer.Addr)
	go func() {
		err := h.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			hsys.Error("# Hest [", h.code, "] listen failed")
			hsys.Exit(errors.New("Hest startup failed: " + err.Error()))
			return
		}
	}()
	return nil, nil
}

func (h *Hest) shutdown(ctx context.Context) {
	if err := h.httpServer.Shutdown(ctx); err != nil {
		hsys.Error("# Rest [", h.code, "] shutdown error: ", err.Error())
		return
	}
	hsys.Error("# Hest [", h.code, "] exist: ", h.httpServer.Addr)
}
