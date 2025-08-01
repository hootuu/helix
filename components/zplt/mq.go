package zplt

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/helix/unicom/hmq/hnsq"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
	"sync"
)

const (
	helixMainMQ         = "helix_main_nsq"
	helixMainMqProducer = "helix_main_nsq_producer"
)

var gMQ *hmq.MQ
var gMqProducer *hmq.Producer
var gMqMu sync.Mutex

func HelixMQ() *hmq.MQ {
	return gMQ
}

func HelixMqProducer() (*hmq.Producer, error) {
	if gMqProducer != nil {
		return gMqProducer, nil
	}
	gMqMu.Lock()
	defer gMqMu.Unlock()
	if gMqProducer != nil {
		return gMqProducer, nil
	}
	producer := HelixMQ().NewProducer()
	err := HelixMQ().RegisterProducer(producer)
	if err != nil {
		return nil, err
	}
	gMqProducer = producer
	return gMqProducer, nil
}

func HelixMqPublish(topic hmq.Topic, payload hmq.Payload) error {
	producer, err := HelixMqProducer()
	if err != nil {
		return err
	}
	return producer.Publish(topic, payload)
}

func HelixMqMustPublish(topic hmq.Topic, payload hmq.Payload) {
	err := hretry.Must(func() error {
		return HelixMqPublish(topic, payload)
	})
	if err != nil {
		hlog.Fix("MQ Publish Error: ",
			zap.Error(err),
			zap.Any("topic", topic),
			zap.ByteString("payload", payload),
		)
	}
}

func init() {
	helix.AfterStartup(func() {
		mqRunning := hcfg.GetBool("helix.mq.running", true)
		if !mqRunning {
			gMQ = hmq.NewMQ(helixMainMQ, hmq.NewEmptyMQ())
			return
		}
		gMQ = hmq.NewMQ(helixMainMQ, hnsq.NewNsqMQ())
		gMqProducer = gMQ.NewProducer()
		err := gMQ.RegisterProducer(gMqProducer)
		if err != nil {
			hsys.Exit(err)
		}
	})
}
