package hnsq

import (
	"context"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/hlog"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type Consumer struct {
	nsq      *NsqMQ
	consumer *nsq.Consumer
	self     *hmq.Consumer
}

func newConsumer(nsq *NsqMQ) *Consumer {
	return &Consumer{nsq: nsq}
}

func (c *Consumer) Startup(self *hmq.Consumer) (context.Context, error) {
	c.self = self
	var err error
	config := nsq.NewConfig()
	c.consumer, err = nsq.NewConsumer(string(self.Topic()), string(self.Channel()), config)
	if err != nil {
		hlog.Err("helix.hnsq.Consumer.Startup", zap.Error(err))
		return nil, err
	}
	c.consumer.AddHandler(c)
	if err := c.consumer.ConnectToNSQD(c.nsq.nsqdAddr); err != nil {
		hlog.Err("helix.hnsq.Consumer.Startup:Connect", zap.Error(err))
		return nil, err
	}
	return nil, nil
}

func (c *Consumer) Shutdown(_ context.Context) {
	if c.consumer != nil {
		c.consumer.Stop()
	}
}

func (c *Consumer) HandleMessage(message *nsq.Message) error {
	return c.self.Handle(
		hmq.NewMessage(string(message.ID[:]), message.Timestamp, message.Body),
	)
}
