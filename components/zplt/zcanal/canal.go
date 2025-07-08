package zcanal

import (
	"github.com/hootuu/helix/storage/hcanal"
)

const (
	helixCanal = "canal"
)

var gHelixCanal *hcanal.Canal

func HelixCanal() *hcanal.Canal {
	return gHelixCanal
}

func init() {
	gHelixCanal = hcanal.New(helixCanal)
}
