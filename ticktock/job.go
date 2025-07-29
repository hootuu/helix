package ticktock

import (
	"context"
	"github.com/hibiken/asynq"
	"time"
)

type JobID = string
type JobTplID = string
type JobType = string
type Payload = []byte

type JobDefinable interface {
	Validate() error
	GetType() JobType
	GetID() JobID
	GetPayload() Payload
	IsPeriodic() bool
	GetUniqueTTL() time.Duration
	ToAsynqTask() *asynq.Task
}

type Job struct {
	Type    JobType `json:"type"`
	Payload Payload `json:"payload"`
}

type JobHandler interface {
	Handle(ctx context.Context, job *Job) error
}

type JobHandlerFunc func(ctx context.Context, job *Job) error
