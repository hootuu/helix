package ticktock

import (
	"errors"
	"github.com/hibiken/asynq"
	"github.com/hootuu/hyle/data/idx"
	"time"
)

type RunAtJob struct {
	Type      JobType       `json:"type"`
	ID        JobID         `json:"id"`
	Payload   Payload       `json:"payload"`
	UniqueTTL time.Duration `json:"unique_ttl"`
	RunAt     time.Time     `json:"run_at"`
}

func (j *RunAtJob) Validate() error {
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

func (j *RunAtJob) GetType() JobType {
	return j.Type
}

func (j *RunAtJob) GetID() JobID {
	return j.ID
}

func (j *RunAtJob) GetPayload() Payload {
	return j.Payload
}

func (j *RunAtJob) GetUniqueTTL() time.Duration {
	return j.UniqueTTL
}

func (j *RunAtJob) IsPeriodic() bool {
	return false
}

func (j *RunAtJob) ToAsynqTask() *asynq.Task {
	opt := []asynq.Option{
		asynq.TaskID(j.ID),
		asynq.ProcessAt(j.RunAt),
		asynq.Queue(qCritical),
	}
	if j.UniqueTTL > 0 {
		opt = append(opt, asynq.Unique(j.UniqueTTL))
	}
	return asynq.NewTask(j.Type, j.Payload, opt...)
}
