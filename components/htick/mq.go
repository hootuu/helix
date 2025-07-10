package htick

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
)

var gTickProducer *hmq.Producer

func triggerTick(job *Job) error {
	err := gTickProducer.Publish(job.Topic, job.Payload)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	helix.AfterStartup(func() {
		gTickProducer = zplt.HelixMQ().NewProducer()
		err := zplt.HelixMQ().RegisterProducer(gTickProducer)
		if err != nil {
			gLogger.Error("htick.init", zap.Error(err))
			hsys.Exit(err)
			return
		}
	})
}
