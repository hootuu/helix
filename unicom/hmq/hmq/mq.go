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
	NewProducer() ProducerCore
	NewConsumer() ConsumerCore
}

type EmptyCore struct {
}

func NewEmptyMQ() Core {
	return &EmptyCore{}
}

func (e EmptyCore) Startup(_ *MQ) (context.Context, error) {
	return nil, nil
}

func (e EmptyCore) Shutdown(_ context.Context) {
}

func (e EmptyCore) NewProducer() ProducerCore {
	return &EmptyProducerCore{}
}

func (e EmptyCore) NewConsumer() ConsumerCore {
	return &EmptyConsumerCore{}
}

type MQ struct {
	Code        string
	core        Core
	consumerArr []*Consumer
	producerArr []*Producer
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
		if len(mq.consumerArr) > 0 {
			for _, consumer := range mq.consumerArr {
				consumer.shutdown()
			}
		}
		mq.core.Shutdown(ctx)
	}
}

func (mq *MQ) NewProducer() *Producer {
	return newProducer(mq.core.NewProducer())
}

func (mq *MQ) NewConsumer(code string, topic Topic, channel Channel) *Consumer {
	return newConsumer(code, topic, channel, mq.core.NewConsumer())
}

func (mq *MQ) RegisterProducer(p *Producer) error {
	err := p.startup()
	if err != nil {
		hlog.Err("hmq.RegisterConsumer", zap.Error(err))
		return err
	}
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.producerArr = append(mq.producerArr, p)
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
	mq.consumerArr = append(mq.consumerArr, c)
	return nil
}
