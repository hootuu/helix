package hmq

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ConsumerCore interface {
	Startup(self *Consumer) (context.Context, error)
	Shutdown(ctx context.Context)
}

type Consumer struct {
	code        string
	topic       Topic
	channel     Channel
	core        ConsumerCore
	coreCtx     context.Context
	handlerFunc func(msg *Message) error
	ctx         context.Context
	ctxStop     context.CancelFunc
	wg          sync.WaitGroup
}

func NewConsumer(code string, topic Topic, channel Channel) *Consumer {
	return &Consumer{
		code:    code,
		topic:   topic,
		channel: channel,
	}
}

func (c *Consumer) Code() string {
	return c.code
}

func (c *Consumer) Topic() Topic {
	return c.topic
}

func (c *Consumer) Channel() Channel {
	return c.channel
}

func (c *Consumer) Handle(msg *Message) error {
	gMqCLogger.Info(c.code, zap.String("code", c.code), zap.String("id", msg.ID),
		zap.String("topic", string(c.topic)), zap.String("channel", string(c.channel)))
	start := time.Now()
	defer func() {
		gMqCLogger.Info(c.code, zap.Int64("_elapse", time.Since(start).Milliseconds()))
	}()
	return c.handlerFunc(msg)
}

func (c *Consumer) With(core ConsumerCore) *Consumer {
	c.core = core
	return c
}

func (c *Consumer) WithHandler(handlerFunc func(msg *Message) error) *Consumer {
	c.handlerFunc = handlerFunc
	return c
}

func (c *Consumer) startup() error {
	if c.core == nil {
		return errors.New("must set consumer core")
	}
	if c.handlerFunc == nil {
		return errors.New("must set handler function")
	}
	c.ctx, c.ctxStop = context.WithCancel(context.Background())
	var err error
	c.coreCtx, err = c.core.Startup(c)
	if err != nil {
		return err
	}
	c.wg.Add(1)
	go c.monitorShutdown()
	return nil
}

func (c *Consumer) shutdown() {
	if c.ctxStop != nil {
		c.ctxStop()
	}
}

func (c *Consumer) monitorShutdown() {
	defer c.wg.Done()

	select {
	case <-c.ctx.Done():
		c.core.Shutdown(c.coreCtx)
	}
}
