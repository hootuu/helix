package hmq

import (
	"errors"
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
	gMqPLogger.Info("publish", zap.String("topic", string(topic)))
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
