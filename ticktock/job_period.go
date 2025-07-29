package ticktock

import (
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"time"
)

type PeriodicJobPayload struct {
	Type         JobType       `json:"type"`
	Payload      Payload       `json:"payload"`
	UniqueTTL    time.Duration `json:"unique_ttl"`
	JobTplID     JobTplID      `json:"job_tpl_id"`
	TplUniqueTTL time.Duration `json:"tpl_unique_ttl"`
	Sequence     uint32        `json:"sequence"`
	Expression   Expression    `json:"expression"`
	Current      time.Time     `json:"current"`
}

type PeriodicJobDefinable interface {
	GetExpression() Expression
	GetCurrent() time.Time
	BuildID() JobID
	BuildPayload() *PeriodicJobPayload
}

type PeriodicJob struct {
	Type         JobType       `json:"type"`
	Payload      Payload       `json:"payload"`
	UniqueTTL    time.Duration `json:"unique_ttl"`
	JobTplID     JobTplID      `json:"job_tpl_id"`
	TplUniqueTTL time.Duration `json:"tpl_unique_ttl"`
	Sequence     uint32        `json:"sequence"`
	Expression   Expression    `json:"expression"`
	Current      time.Time     `json:"current"`
}

func PeriodicJobFromPayload(payload *PeriodicJobPayload) *PeriodicJob {
	return &PeriodicJob{
		Type:         payload.Type,
		Payload:      payload.Payload,
		UniqueTTL:    payload.UniqueTTL,
		JobTplID:     payload.JobTplID,
		TplUniqueTTL: payload.TplUniqueTTL,
		Sequence:     payload.Sequence,
		Expression:   payload.Expression,
		Current:      payload.Current,
	}
}

func (j *PeriodicJob) Validate() error {
	if j.Type == "" {
		return errors.New("type is required")
	}
	if j.JobTplID == "" {
		return errors.New("job_tpl_id is required")
	}
	if len(j.Payload) == 0 {
		return errors.New("payload is required")
	}
	if err := j.Expression.Validate(); err != nil {
		return err
	}
	return nil
}

func (j *PeriodicJob) GetType() JobType {
	return j.Type
}

func (j *PeriodicJob) GetID() JobID {
	return j.BuildID()
}

func (j *PeriodicJob) BuildID() JobID {
	return fmt.Sprintf("%s:%d:%s", j.JobTplID, j.Sequence, idx.New())
}

func (j *PeriodicJob) GetPayload() Payload {
	return j.Payload
}

func (j *PeriodicJob) GetUniqueTTL() time.Duration {
	return j.UniqueTTL
}

func (j *PeriodicJob) GetJobTplID() JobTplID {
	return j.JobTplID
}

func (j *PeriodicJob) GetTplUniqueTTL() time.Duration {
	return j.TplUniqueTTL
}

func (j *PeriodicJob) GetSequence() uint32 {
	return j.Sequence
}

func (j *PeriodicJob) GetExpression() Expression {
	return j.Expression
}

func (j *PeriodicJob) GetCurrent() time.Time {
	return j.Current
}

func (j *PeriodicJob) ToAsynqTask() *asynq.Task {
	nxtTime, err := j.Expression.Next(j.Current)
	if err != nil {
		hlog.Fix("should not catch this line, because validate in NewPeriodicJob",
			zap.String("jobType", j.Type),
			zap.String("jobTplID", j.JobTplID),
			zap.Uint32("sequence", j.Sequence),
			zap.Error(err))
		nxtTime = time.Now()
	}
	opt := []asynq.Option{
		asynq.TaskID(j.GetID()),
		asynq.ProcessAt(nxtTime),
	}
	if j.UniqueTTL > 0 {
		opt = append(opt, asynq.Unique(j.UniqueTTL))
	}
	return asynq.NewTask(j.Type, j.Payload, opt...)
}

func (j *PeriodicJob) IsPeriodic() bool {
	return true
}

func (j *PeriodicJob) BuildPayload() *PeriodicJobPayload {
	nxtTime, err := j.Expression.Next(j.Current)
	if err != nil {
		hlog.Fix("should not catch this line, because validate has been call",
			zap.String("jobType", j.Type),
			zap.String("jobTplID", j.JobTplID),
			zap.Uint32("sequence", j.Sequence),
			zap.Error(err))
		nxtTime = time.Now()
	}
	//nxtTime, _ = j.Expression.Next(nxtTime)
	fmt.Println("j.Current:", j.Current.Format(time.RFC3339))
	fmt.Println("j.NxtTime:", nxtTime.Format(time.RFC3339))
	return &PeriodicJobPayload{
		Type:         j.Type,
		Payload:      j.Payload,
		UniqueTTL:    j.UniqueTTL,
		JobTplID:     j.JobTplID,
		TplUniqueTTL: j.TplUniqueTTL,
		Sequence:     j.Sequence + 1,
		Expression:   j.Expression,
		Current:      nxtTime,
	}
}
