package hmq

import (
	"context"
	"errors"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"sync"
)

type Core interface {
	Startup(self *MQ) (context.Context, error)
	Shutdown(ctx context.Context)
}

type MQ struct {
	Code        string
	core        Core
	producerMap map[string]*Producer
	consumerMap map[string]*Consumer
	mu          sync.Mutex
}

func NewMQ(code string, core Core) *MQ {
	mq := &MQ{Code: code, core: core}
	helix.Use(mq.Helix())
	return mq
}

func (mq *MQ) Helix() helix.Helix {
	return helix.BuildHelix(mq.Code, mq.startup, mq.shutdown)
}

func (mq *MQ) startup() (context.Context, error) {
	if mq.core == nil {
		return nil, errors.New("must set consumer core")
	}
	return mq.core.Startup(mq)
}

func (mq *MQ) shutdown(ctx context.Context) {
	if mq.core != nil {
		if len(mq.producerMap) > 0 {
			for _, producer := range mq.producerMap {
				producer.shutdown()
			}
			for _, consumer := range mq.consumerMap {
				consumer.shutdown()
			}
		}
		mq.core.Shutdown(ctx)
	}
}

func (mq *MQ) RegisterProducer(p *Producer) error {
	err := p.startup()
	if err != nil {
		hlog.Err("hmq.RegisterProducer", zap.Error(err))
		return err
	}
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if mq.producerMap == nil {
		mq.producerMap = make(map[string]*Producer)
	}
	if _, ok := mq.producerMap[p.code]; !ok {
		mq.producerMap[p.code] = p
	}
	return nil
}

func (mq *MQ) RegisterConsumer(c *Consumer) error {
	err := c.startup()
	if err != nil {
		hlog.Err("hmq.RegisterConsumer", zap.Error(err))
		return err
	}
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if mq.consumerMap == nil {
		mq.consumerMap = make(map[string]*Consumer)
	}
	if _, ok := mq.consumerMap[c.code]; !ok {
		mq.consumerMap[c.code] = c
	}
	return nil
}
