package hnsq

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/hcfg"
)

type NsqMQ struct {
	self     *hmq.MQ
	nsqdAddr string
}

func NewNsqMQ() *NsqMQ {
	return &NsqMQ{}
}

func (n *NsqMQ) NewProducer() hmq.ProducerCore {
	return newProducer(n)
}

func (n *NsqMQ) NewConsumer() hmq.ConsumerCore {
	return newConsumer(n)
}

func (n *NsqMQ) Startup(self *hmq.MQ) (context.Context, error) {
	n.self = self
	n.nsqdAddr = hcfg.GetString(n.cfgKey("nsqd.addr"), "127.0.0.1:4150")
	return nil, nil
}

func (n *NsqMQ) Shutdown(_ context.Context) {

}

func (n *NsqMQ) cfgKey(key string) string {
	return fmt.Sprintf("helix.hnsq.%s.%s", n.self.Code, key)
}
