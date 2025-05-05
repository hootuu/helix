package zplt

import (
	"github.com/hootuu/helix/storage/hrds"
)

const (
	helixRdsCache = "helix"
)

func HelixRdsCache() *hrds.Cache {
	return hrds.GetCache(helixRdsCache)
}

func init() {
	hrds.Register(helixRdsCache)
}
