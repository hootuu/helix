package zplt

import (
	"github.com/hootuu/helix/storage/hrds"
)

const (
	helixRdsCache = "helix_rds"
)

func HelixRdsCache() *hrds.Cache {
	return hrds.GetCache(helixRdsCache)
}

func init() {
	hrds.Register(helixRdsCache)
}
