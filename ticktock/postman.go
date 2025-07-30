package ticktock

import (
	"context"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hrds"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"time"
)

type Postman struct {
	code  string
	cache *hrds.Cache
	cli   *asynq.Client
}

func NewPostman(code string, cache *hrds.Cache) *Postman {
	postman := &Postman{
		code:  code,
		cache: cache,
		cli:   asynq.NewClientFromRedisClient(cache.Redis()),
	}
	helix.Use(postman.Helix())
	return postman
}

func (p *Postman) Helix() helix.Helix {
	return helix.BuildHelix(
		p.code,
		p.doStartup,
		p.doShutdown,
	)
}

func (p *Postman) Code() string {
	return p.code
}

func (p *Postman) Send(ctx context.Context, job JobDefinable) (err error) {
	if job == nil {
		return errors.New("job is nil")
	}
	if err = job.Validate(); err != nil {
		return err
	}
	var postID string
	defer func() {
		if err != nil {
			gLogger.Error("[OUT]"+job.GetType(),
				hlog.TraceInfo(ctx), zap.Error(err),
				zap.String("job_type", job.GetType()),
				zap.String("job_id", job.GetID()),
				zap.String("payload", string(job.GetPayload())))
		} else {
			if hlog.IsElapseDetail() {
				gLogger.Info("[OUT]"+job.GetType(), hlog.TraceInfo(ctx),
					zap.String("job_type", job.GetType()),
					zap.String("job_id", job.GetID()),
					zap.String("post_id", postID))
			}
		}
	}()
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "ticktock.Send",
			hlog.F(zap.String("job_type", job.GetType()),
				zap.String("job_id", job.GetID())),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err), zap.String("job_payload", string(job.GetPayload()))}
				}
				return []zap.Field{zap.String("post_id", postID)}
			},
		)()
	}
	if job.IsPeriodic() {
		periodicJob, ok := job.(*PeriodicJob)
		if !ok {
			return errors.New("job does not implement PeriodicJobDefinable when it is Periodic")
		}
		periodicPayload := periodicJob.BuildPayload()
		periodExpr := periodicJob.GetExpression()
		nextTime, err := periodExpr.Next(periodicJob.GetCurrent())
		if err != nil {
			return err
		}
		nxtPeriodJob := &RunAtJob{
			Type:      periodicTaskType,
			ID:        periodicJob.JobTplID + fmt.Sprintf("_%d", periodicJob.Sequence),
			Payload:   hjson.MustToBytes(periodicPayload),
			UniqueTTL: periodicJob.TplUniqueTTL,
			RunAt:     nextTime.Add(-10 * time.Second),
		}
		_, err = p.doSend(ctx, nxtPeriodJob)
		if err != nil {
			return err
		}
		if job.GetType() == MqTaskType {
			hlog.Info(fmt.Sprintf("nextTime: %s", nextTime.Format("2006-01-02 15:04:05")),
				hlog.TraceInfo(ctx), zap.String("job_type", job.GetType()))
			LocalSchedule(nextTime, func() {
				err = onMqJobHandlerFunc(ctx, &Job{
					Type:    MqTaskType,
					Payload: job.GetPayload(),
				})
			})
		} else {
			itemDoJob := &RunAtJob{
				Type:      periodicJob.GetType(),
				ID:        periodicJob.BuildID(),
				Payload:   periodicJob.GetPayload(),
				UniqueTTL: periodicJob.UniqueTTL,
				RunAt:     nextTime,
			}
			postID, err = p.doSend(ctx, itemDoJob)
			if err != nil {
				return err
			}
		}
		return nil
	}
	postID, err = p.doSend(ctx, job)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postman) doSend(ctx context.Context, job JobDefinable) (id string, err error) {
	task := job.ToAsynqTask()
	info, err := p.cli.Enqueue(task)
	if err != nil {
		hlog.TraceErr("asynq.cli.Enqueue failed", ctx, err)
		return "", err
	}
	return info.ID, nil
}

func (p *Postman) doStartup() (context.Context, error) {
	return nil, nil
}

func (p *Postman) doShutdown(_ context.Context) {
	if p.cli != nil {
		err := p.cli.Close()
		if err != nil {
			hlog.Err("ticktock.doShutdown", zap.String("code", p.code), zap.Error(err))
		}
	}
}
