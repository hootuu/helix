package zmeili

import (
	"github.com/hootuu/helix/storage/hmeili"
)

const (
	helixMeili = "meili"
)

func HelixMeili() *hmeili.Meili {
	return hmeili.GetMeili(helixMeili)
}

func init() {
	hmeili.Register(helixMeili)
}
