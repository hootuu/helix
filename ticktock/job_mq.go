package ticktock

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
)

const (
	MqTaskType = "TICKTOCK_MQ"
)

type MqJobPayload struct {
	Topic   hmq.Topic   `json:"topic"`
	Payload hmq.Payload `json:"payload"`
}

func onMqJobHandlerFunc(ctx context.Context, job *Job) error {
	if job == nil {
		hlog.TraceFix("ticktock: the job should not be nil", ctx, fmt.Errorf("nil job"))
		return nil
	}
	if len(job.Payload) == 0 {
		hlog.TraceFix("ticktock: the job.payload is empty", ctx, fmt.Errorf("nil job"))
		return nil
	}
	payload := hjson.MustFromBytes[MqJobPayload](job.Payload)
	if payload == nil {
		hlog.TraceFix("ticktock: the mq.payload is empty", ctx, fmt.Errorf("empty mq payload"))
		return nil
	}
	zplt.HelixMqMustPublish(payload.Topic, payload.Payload)
	return nil
}
