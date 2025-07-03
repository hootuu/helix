package hnsq

import (
	"context"
	"errors"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/hlog"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type Producer struct {
	nsq      *NsqMQ
	producer *nsq.Producer
	self     *hmq.Producer
}

func newProducer(nsq *NsqMQ) *Producer {
	return &Producer{nsq: nsq}
}

func (p *Producer) Startup(self *hmq.Producer) (context.Context, error) {
	p.self = self
	var err error
	config := nsq.NewConfig()
	p.producer, err = nsq.NewProducer(p.nsq.nsqdAddr, config)
	if err != nil {
		hlog.Err("helix.hnsq.Producer.Startup", zap.Error(err))
		return nil, err
	}
	if err := p.producer.Ping(); err != nil {
		hlog.Err("helix.hnsq.Producer.Startup:Ping", zap.Error(err))
		return nil, err
	}

	return nil, nil
}

func (p *Producer) Shutdown(_ context.Context) {
	if p.producer != nil {
		p.producer.Stop()
	}
}

func (p *Producer) Publish(payload hmq.Payload) error {
	if p.producer == nil {
		return errors.New("nsq startup failed")
	}
	return p.producer.Publish(string(p.self.Topic()), payload)
}
