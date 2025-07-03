package hmq

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

type ProducerCore interface {
	Startup(self *Producer) (context.Context, error)
	Shutdown(ctx context.Context)
	Publish(payload Payload) error
}

type Producer struct {
	code    string
	topic   Topic
	core    ProducerCore
	coreCtx context.Context
}

func NewProducer(code string, topic Topic) *Producer {
	return &Producer{
		code:  code,
		topic: topic,
	}
}

func (p *Producer) Code() string {
	return p.code
}

func (p *Producer) Topic() Topic {
	return p.topic
}

func (p *Producer) Publish(payload Payload) error {
	gMqPLogger.Info(p.code, zap.String("code", p.code), zap.String("topic", string(p.topic)))
	if p.core == nil {
		return errors.New("must set producer core")
	}
	if payload == nil {
		return errors.New("must set payload")
	}
	return p.core.Publish(payload)
}

func (p *Producer) With(core ProducerCore) *Producer {
	p.core = core
	return p
}

func (p *Producer) startup() error {
	if p.core == nil {
		return errors.New("must set producer core")
	}
	var err error
	p.coreCtx, err = p.core.Startup(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) shutdown() {
	p.core.Shutdown(p.coreCtx)
}
