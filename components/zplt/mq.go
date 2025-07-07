package zplt

import (
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/helix/unicom/hmq/hnsq"
)

const (
	helixMainMQ = "helix_main_nsq"
)

var gMQ *hmq.MQ

func HelixMQ() *hmq.MQ {
	return gMQ
}

func init() {
	gMQ = hmq.NewMQ(helixMainMQ, hnsq.NewNsqMQ())
}
