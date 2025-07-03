package zplt

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"gorm.io/gorm"
)

const (
	helixDb = "helix_mysql"
)

func HelixPgDB() *hdb.Database {
	return hdb.GetDatabase(helixDb)
}

func HelixPgCtx(ctx context.Context) *gorm.DB {
	tx := hdb.CtxTx(ctx)
	if tx == nil {
		tx = HelixPgDB().PG()
	}
	return tx
}

func HelixPgTx(ctx ...context.Context) context.Context {
	return hdb.TxCtx(HelixPgDB().PG(), ctx...)
}

func init() {
	hdb.Register(helixDb)
}
