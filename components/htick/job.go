package htick

import (
	"errors"
	"github.com/hootuu/helix/unicom/hmq/hmq"
)

type Job struct {
	Expression Expression  `json:"expression"`
	Topic      hmq.Topic   `json:"topic"`
	Payload    hmq.Payload `json:"payload"`
}

func NewJob(exp Expression, topic hmq.Topic, p hmq.Payload) *Job {
	return &Job{
		Expression: exp,
		Topic:      topic,
		Payload:    p,
	}
}

func (job *Job) Validate() error {
	if job.Expression == "" {
		return errors.New("job expression is required")
	}
	if job.Topic == "" {
		return errors.New("job topic is required")
	}
	if len(job.Payload) == 0 {
		return errors.New("job payload is required")
	}
	return nil
}
