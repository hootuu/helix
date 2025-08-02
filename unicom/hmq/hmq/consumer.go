package hmq

import (
	"context"
	"errors"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ConsumerCore interface {
	Startup(self *Consumer) (context.Context, error)
	Shutdown(ctx context.Context)
}

type EmptyConsumerCore struct{}

func (e *EmptyConsumerCore) Startup(_ *Consumer) (context.Context, error) {
	return nil, nil
}

func (e *EmptyConsumerCore) Shutdown(_ context.Context) {

}

type Consumer struct {
	code        string
	topic       Topic
	channel     Channel
	core        ConsumerCore
	coreCtx     context.Context
	handlerFunc func(ctx context.Context, msg *Message) error
	ctx         context.Context
	ctxStop     context.CancelFunc
	wg          sync.WaitGroup
}

func newConsumer(code string, topic Topic, channel Channel, core ConsumerCore) *Consumer {
	return &Consumer{
		code:        code,
		topic:       topic,
		channel:     channel,
		core:        core,
		handlerFunc: func(ctx context.Context, msg *Message) error { return nil },
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
	var err error
	ctx := context.WithValue(context.Background(), hlog.TraceIdKey, msg.ID)
	if hlog.IsElapseDetail() {
		gMqCLogger.Info(msg.ID, zap.String("code", c.code),
			hlog.TraceInfo(ctx),
			zap.String("id", msg.ID),
			zap.String("topic", string(c.topic)),
			zap.String("channel", string(c.channel)),
			zap.String("payload", string(msg.Payload)))
		start := time.Now()
		defer func() {
			arr := []zap.Field{
				zap.Int64("_elapse", time.Since(start).Milliseconds()),
				hlog.TraceInfo(ctx),
			}
			if err != nil {
				arr = append(arr, zap.Error(err), zap.String("payload", string(msg.Payload)))
				gMqCLogger.Error(msg.ID, arr...)
				return
			}
			gMqCLogger.Info(msg.ID, arr...)
		}()
	}
	err = c.handlerFunc(ctx, msg)
	if err != nil {
		hlog.Err("hmq.consumer.Handle",
			hlog.TraceInfo(ctx),
			zap.String("code", c.code),
			zap.String("id", msg.ID),
			zap.String("topic", string(c.topic)),
			zap.String("channel", string(c.channel)),
			zap.String("payload", string(msg.Payload)),
			zap.Error(err))
		return err
	}
	return nil
}

func (c *Consumer) WithHandler(handlerFunc func(ctx context.Context, msg *Message) error) *Consumer {
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
