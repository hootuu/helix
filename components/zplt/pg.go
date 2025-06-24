package zplt

import (
	"context"
	"github.com/hootuu/helix/storage/hpg"
	"gorm.io/gorm"
)

const (
	helixDb = "helix_pg"
)

func HelixPgDB() *hpg.Database {
	return hpg.GetDatabase(helixDb)
}

func HelixPgCtx(ctx context.Context) *gorm.DB {
	tx := hpg.CtxTx(ctx)
	if tx == nil {
		tx = HelixPgDB().PG()
	}
	return tx
}

func HelixPgTx(ctx ...context.Context) context.Context {
	return hpg.TxCtx(HelixPgDB().PG(), ctx...)
}

func init() {
	hpg.Register(helixDb)
}
