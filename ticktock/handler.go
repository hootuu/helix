package ticktock

import (
	"context"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func jobFromAsynqTask(task *asynq.Task) (*Job, error) {
	if task == nil {
		return nil, errors.New("task is nil")
	}
	return &Job{
		Type:    task.Type(),
		Payload: task.Payload(),
	}, nil
}

type asyncHandlerWrapper struct {
	handler JobHandler
}

func newAsyncHandlerWrapper(handler JobHandler) *asyncHandlerWrapper {
	return &asyncHandlerWrapper{handler: handler}
}

func (w *asyncHandlerWrapper) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	job, err := jobFromAsynqTask(task)
	if err != nil {
		return err
	}
	innerCtx := hlog.NewTraceCtx(ctx)
	defer func() {
		if err != nil {
			gLogger.Error("[IN]"+job.Type, hlog.TraceInfo(ctx), zap.Error(err), zap.String("payload", string(job.Payload)))
		} else {
			if hlog.IsElapseDetail() {
				gLogger.Info("[IN]"+job.Type, hlog.TraceInfo(ctx), zap.String("payload", string(job.Payload)))
			}
		}
	}()
	return w.handler.Handle(innerCtx, job)
}

func asynqHandleFuncWrapper(handleFunc JobHandlerFunc) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) (err error) {
		if task == nil {
			hlog.Fix("ticktock: the ticktock task should not be nil")
			return nil
		}
		payload := task.Payload()
		if len(payload) == 0 {
			hlog.Fix("ticktock: the len of ticktock task.payload should not be 0")
			return nil
		}
		innerCtx := hlog.NewTraceCtx(ctx)
		job, err := jobFromAsynqTask(task)

		defer func() {
			if err != nil {
				gLogger.Error("[IN]"+job.Type, hlog.TraceInfo(innerCtx), zap.Error(err), zap.String("payload", string(job.Payload)))
			} else {
				if hlog.IsElapseDetail() {
					gLogger.Info("[IN]"+job.Type, hlog.TraceInfo(innerCtx), zap.String("payload", string(job.Payload)))
				}
			}
		}()
		return handleFunc(innerCtx, job)
	}
}
