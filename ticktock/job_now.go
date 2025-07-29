package ticktock

import (
	"errors"
	"github.com/hibiken/asynq"
	"github.com/hootuu/hyle/data/idx"
	"time"
)

type NowJob struct {
	Type      JobType       `json:"type"`
	ID        JobID         `json:"id"`
	Payload   Payload       `json:"payload"`
	UniqueTTL time.Duration `json:"unique_ttl"`
}

func (j *NowJob) Validate() error {
	if j.Type == "" {
		return errors.New("type is required")
	}
	if j.ID == "" {
		j.ID = j.Type + ":" + idx.New()
	}
	if len(j.Payload) == 0 {
		return errors.New("payload is required")
	}
	return nil
}

func (j *NowJob) GetType() JobType {
	return j.Type
}

func (j *NowJob) GetID() JobID {
	return j.ID
}

func (j *NowJob) GetPayload() Payload {
	return j.Payload
}

func (j *NowJob) GetUniqueTTL() time.Duration {
	return j.UniqueTTL
}

func (j *NowJob) IsPeriodic() bool {
	return false
}

func (j *NowJob) ToAsynqTask() *asynq.Task {
	opt := []asynq.Option{asynq.TaskID(j.ID)}
	if j.UniqueTTL > 0 {
		opt = append(opt, asynq.Unique(j.UniqueTTL))
	}
	return asynq.NewTask(j.Type, j.Payload, opt...)
}
