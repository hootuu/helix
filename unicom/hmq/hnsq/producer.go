package hnsq

import (
	"errors"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/hlog"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type Producer struct {
	nsq      *NsqMQ
	producer *nsq.Producer
}

func newProducer(nsq *NsqMQ) *Producer {
	return &Producer{nsq: nsq}
}

func (p *Producer) Startup() error {
	var err error
	config := nsq.NewConfig()
	p.producer, err = nsq.NewProducer(p.nsq.nsqdAddr, config)
	if err != nil {
		hlog.Err("helix.hnsq.Producer.Startup", zap.Error(err))
		return err
	}
	if err := p.producer.Ping(); err != nil {
		hlog.Err("helix.hnsq.Producer.Startup:Ping", zap.Error(err))
		return err
	}

	return nil
}

func (p *Producer) Shutdown() {
	if p.producer != nil {
		p.producer.Stop()
	}
}

func (p *Producer) Publish(topic hmq.Topic, payload hmq.Payload) error {
	if p.producer == nil {
		return errors.New("nsq startup failed")
	}
	return p.producer.Publish(string(topic), payload)
}
