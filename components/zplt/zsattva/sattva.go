package zsattva

import (
	"github.com/hootuu/helix/components/sattva"
)

const (
	helixMainSattva = "helix_main_sattva"
)

var gMainSattva *sattva.Sattva

func HelixSattva() *sattva.Sattva {
	return gMainSattva
}

func init() {
	gMainSattva, _ = sattva.NewSattva(helixMainSattva)
}
