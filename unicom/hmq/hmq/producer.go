package hmq

import (
	"errors"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

type ProducerCore interface {
	Startup() error
	Shutdown()
	Publish(topic Topic, payload Payload) error
}

type Producer struct {
	core ProducerCore
}

func newProducer(core ProducerCore) *Producer {
	return &Producer{
		core: core,
	}
}

func (p *Producer) Publish(topic Topic, payload Payload) error {
	if hlog.IsElapseDetail() {
		gMqPLogger.Info("publish",
			zap.String("topic", string(topic)),
			zap.String("payload", string(payload)),
		)
	}
	if p.core == nil {
		return errors.New("must set producer core")
	}
	return p.core.Publish(topic, payload)
}

func (p *Producer) startup() error {
	if p.core == nil {
		return errors.New("must set producer core")
	}
	var err error
	err = p.core.Startup()
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) shutdown() {
	p.core.Shutdown()
}
