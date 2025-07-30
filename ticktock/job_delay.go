package ticktock

import (
	"errors"
	"github.com/hibiken/asynq"
	"github.com/hootuu/hyle/data/idx"
	"time"
)

type DelayJob struct {
	Type      JobType       `json:"type"`
	ID        JobID         `json:"id"`
	Payload   Payload       `json:"payload"`
	UniqueTTL time.Duration `json:"unique_ttl"`
	Delay     time.Duration `json:"delay"`
}

func (j *DelayJob) Validate() error {
	if j.Type == "" {
		return errors.New("type is required")
	}
	if j.ID == "" {
		j.ID = j.Type + ":" + idx.New()
	}
	if len(j.Payload) == 0 {
		return errors.New("payload is required")
	}
	if j.Delay <= 0 {
		return errors.New("delay <= 0")
	}
	return nil
}

func (j *DelayJob) GetType() JobType {
	return j.Type
}

func (j *DelayJob) GetID() JobID {
	return j.ID
}

func (j *DelayJob) GetPayload() Payload {
	return j.Payload
}

func (j *DelayJob) GetUniqueTTL() time.Duration {
	return j.UniqueTTL
}

func (j *DelayJob) IsPeriodic() bool {
	return false
}

func (j *DelayJob) ToAsynqTask() *asynq.Task {
	opt := []asynq.Option{
		asynq.TaskID(j.ID),
		asynq.ProcessIn(j.Delay),
		asynq.Queue(qCritical),
	}
	if j.UniqueTTL > 0 {
		opt = append(opt, asynq.Unique(j.UniqueTTL))
	}
	return asynq.NewTask(j.Type, j.Payload, opt...)
}
