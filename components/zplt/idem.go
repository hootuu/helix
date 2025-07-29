package zplt

import (
	"github.com/hootuu/helix/components/hidem"
	"github.com/hootuu/helix/helix"
	"time"
)

const (
	gUniIdemCode = "helix_idem"
)

var gUniIdem hidem.Factory

func UniIdem() hidem.Factory {
	if gUniIdem == nil {
		helix.MustInit(gUniIdemCode, func() error {
			var err error
			gUniIdem, err = hidem.NewCacheFactory(HelixRdsCache(), gUniIdemCode, 15*time.Minute)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return gUniIdem
}
