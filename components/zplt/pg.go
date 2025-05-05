package zplt

import (
	"github.com/hootuu/helix/storage/hpg"
)

const (
	helixDb = "helix"
)

func HelixPgDB() *hpg.Database {
	return hpg.GetDatabase(helixDb)
}

func init() {
	hpg.Register(helixDb)
}
