package hcanal

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"math/rand"
)

type Canal struct {
	Code string
	canal.DummyEventHandler
	alterHandlerArr []AlterHandler
	core            *canal.Canal
}

func New(code string, alterHandler ...AlterHandler) *Canal {
	h := &Canal{Code: code, alterHandlerArr: alterHandler}
	helix.Use(h.Helix())
	return h
}

func (h *Canal) Helix() helix.Helix {
	return helix.BuildHelix(h.Code, h.Startup, h.Shutdown)
}

func (h *Canal) Startup() (context.Context, error) {
	var err error
	cfg := canal.NewDefaultConfig()
	cfg.Addr = hcfg.GetString(h.cfg("addr"), "127.0.0.1:3306")
	cfg.User = hcfg.GetString(h.cfg("user"), "root")
	cfg.Password = hcfg.GetString(h.cfg("password"), "88888888")
	cfg.Flavor = hcfg.GetString(h.cfg("flavor"), "mysql")

	//todo
	cfg.ServerID = uint32(10000 + rand.Intn(10000))
	//todo
	cfg.Dump.ExecutionPath = hcfg.GetString(h.cfg("dump.execution.path"), "/usr/local/mysql/bin/mysqldump")

	//todo
	// We only care table canal_test in test db
	// cfg.Dump.TableDB = "test"
	// cfg.Dump.Tables = []string{"canal_test"}
	//if len(h.alterHandlerArr) > 0 {
	//	for _, alterHandler := range h.alterHandlerArr {
	//		h.core.AddDumpDatabases(alterHandler.Schema())
	//		tables := alterHandler.Table()
	//		h.core.AddDumpTables(alterHandler.Schema(), tables...)
	//	}
	//}

	h.core, err = canal.NewCanal(cfg)
	if err != nil {
		hlog.Err("helix.canal.Startup", zap.String("code", h.Code), zap.Error(err))
		return nil, err
	}

	h.core.SetEventHandler(h)

	go func() {
		err := h.core.Run()
		if err != nil {
			hlog.Err("helix.canal.Startup", zap.String("code", h.Code), zap.Error(err))
		}
	}()

	return nil, nil
}

func (h *Canal) Shutdown(_ context.Context) {
	if h.core != nil {
		h.core.Close()
	}
}

func (h *Canal) cfg(fix string) string {
	return fmt.Sprintf("helix.canal.%s.%s", h.Code, fix)
}
