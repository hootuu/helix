package hpg

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Database struct {
	code string
	pgDB *gorm.DB
}

func newDatabase(code string) *Database {
	return &Database{
		code: code,
	}
}

func (db *Database) Helix() helix.Helix {
	return helix.BuildHelix(db.code, db.startup, db.shutdown)
}

func (db *Database) startup() (context.Context, error) {
	hsys.Info("\n# Connecting to pg.db [", db.code, "] ... #")
	err := retry.Do(
		func() error {
			var err error
			db.pgDB, err = gorm.Open(
				postgres.Open(hcfg.GetString("pg."+db.code+".dns", "host=localhost dbname="+db.code+" port=5432 sslmode=disable")),
				&gorm.Config{
					SkipDefaultTransaction:                   hcfg.GetBool("pg."+db.code+".skip.default.transaction", false),
					NamingStrategy:                           nil,
					FullSaveAssociations:                     hcfg.GetBool("pg."+db.code+".full.save.associations", false),
					Logger:                                   nil,
					NowFunc:                                  nil,
					DryRun:                                   hcfg.GetBool("pg."+db.code+".dry.run", false),
					PrepareStmt:                              hcfg.GetBool("pg."+db.code+".prepare.stmt", false),
					PrepareStmtMaxSize:                       hcfg.GetInt("pg."+db.code+".prepare.stmt.max.size", 0),
					PrepareStmtTTL:                           hcfg.GetDuration("pg."+db.code+".prepare.stmt.ttl", 0),
					DisableAutomaticPing:                     hcfg.GetBool("pg."+db.code+".disable.automatic.ping", false),
					DisableForeignKeyConstraintWhenMigrating: hcfg.GetBool("pg."+db.code+".disable.foreign.key.constraint.when.migrating", false),
					IgnoreRelationshipsWhenMigrating:         hcfg.GetBool("pg."+db.code+".ignore.relationships.when.migrating", false),
					DisableNestedTransaction:                 hcfg.GetBool("pg."+db.code+".disable.nested.transaction", false),
					AllowGlobalUpdate:                        hcfg.GetBool("pg."+db.code+".allow.global.update", false),
					QueryFields:                              hcfg.GetBool("pg."+db.code+".query.fields", false),
					CreateBatchSize:                          hcfg.GetInt("pg."+db.code+".create.batch.size", 0),
					TranslateError:                           hcfg.GetBool("pg."+db.code+".translate.error", false),
					PropagateUnscoped:                        hcfg.GetBool("pg."+db.code+".propagate.unscoped", false),
					ClauseBuilders:                           nil,
					ConnPool:                                 nil,
					Dialector:                                nil,
					Plugins:                                  nil,
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(10),
		retry.Delay(5*time.Second),
	)
	if err != nil {
		hsys.Error("# Connecting to db [", db.code, "] Err:"+err.Error()+"#\n")
		return nil, err
	}
	hsys.Success("# Connecting to db [", db.code, "] OK #\n")
	return context.Background(), nil
}

func (db *Database) shutdown(_ context.Context) {
}

func (db *Database) PG() *gorm.DB {
	return db.pgDB
}

var gPostgresDbMap = make(map[string]*Database)
var gPostgresDbMutex sync.Mutex

func doRegister(code string) {
	gPostgresDbMutex.Lock()
	defer gPostgresDbMutex.Unlock()
	if _, ok := gPostgresDbMap[code]; ok {
		hlog.Err("hpg.doRegister: pg db repetition")
		return
	}
	db := newDatabase(code)
	gPostgresDbMap[code] = db
	helix.Use(db.Helix())
}

func doGetDb(code string) *Database {
	gPostgresDbMutex.Lock()
	defer gPostgresDbMutex.Unlock()
	db, ok := gPostgresDbMap[code]
	if !ok {
		return nil
	}
	return db
}

func doDbActWithRetry(call func() error) error {
	return retry.Do(
		call,
		retry.Attempts(uint(hcfg.GetInt("hpg.act.retry.attempts", 3))),
		retry.Delay(hcfg.GetDuration("hpg.act.retry.delay", 600*time.Millisecond)),
	)
}
