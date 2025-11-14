package ticktock

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hrds"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"runtime"
	"strings"
	"time"
)

type Worker struct {
	code         string
	cache        *hrds.Cache
	srv          *asynq.Server
	srvMux       *asynq.ServeMux
	localPostman *Postman
	baseCtx      context.Context
	stopFunc     context.CancelFunc
}

func NewWorker(code string, cache *hrds.Cache) *Worker {
	w := &Worker{
		code:  code,
		cache: cache,
	}
	w.baseCtx, w.stopFunc = context.WithCancel(context.WithValue(context.Background(), "ticktock", w.code))
	w.srv = asynq.NewServerFromRedisClient(
		w.cache.Redis(),
		asynq.Config{
			Concurrency: hcfg.GetInt(w.cKey("concurrency"), runtime.NumCPU()*4),
			BaseContext: func() context.Context {
				return w.baseCtx
			},
			TaskCheckInterval: hcfg.GetDuration(w.cKey("task.check.interval"), 2*time.Second),
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				switch n {
				case 1:
					return 300 * time.Millisecond
				case 2:
					return 1200 * time.Millisecond
				case 3:
					return 5 * time.Second
				case 4:
					return 10 * time.Second
				}
				return 30 * time.Second
			},
			IsFailure: func(err error) bool {
				return err != nil
			},
			Queues: map[string]int{
				qCritical: 36,
				qDefault:  20,
				qLow:      6,
			},
			StrictPriority:           hcfg.GetBool(w.cKey("strict.priority"), false),
			ErrorHandler:             w,
			Logger:                   newLogger(w.code),
			LogLevel:                 asynq.LogLevel(hcfg.GetInt(w.cKey("log.level"), 2)),
			ShutdownTimeout:          10 * time.Second,
			HealthCheckFunc:          nil,
			HealthCheckInterval:      0,
			DelayedTaskCheckInterval: 0,
			GroupGracePeriod:         0,
			GroupMaxDelay:            0,
			GroupMaxSize:             0,
			GroupAggregator:          nil,
			JanitorInterval:          0,
			JanitorBatchSize:         0,
		},
	)

	w.srvMux = asynq.NewServeMux()

	w.localPostman = NewPostman(w.code+"_cli", cache)

	w.RegisterJobHandlerFunc(periodicTaskType, w.onPeriodicJobHandlerFunc)
	w.RegisterJobHandlerFunc(MqTaskType, onMqJobHandlerFunc)

	helix.Use(w.Helix())

	return w
}

func (w *Worker) Helix() helix.Helix {
	return helix.BuildHelix(w.code, w.doStartup, w.doShutdown)
}

func (w *Worker) GetCode() string {
	return w.code
}

func (w *Worker) RegisterJobHandler(pattern string, handler JobHandler) {
	if w.srvMux == nil {
		return
	}
	w.srvMux.Handle(pattern, newAsyncHandlerWrapper(handler))
}

func (w *Worker) RegisterJobHandlerFunc(pattern string, handlerFunc JobHandlerFunc) {
	if w.srvMux == nil {
		return
	}
	w.srvMux.HandleFunc(pattern, asynqHandleFuncWrapper(handlerFunc))
}

func (w *Worker) onPeriodicJobHandlerFunc(ctx context.Context, job *Job) (err error) {
	if job == nil {
		hlog.Fix("ticktock: the job should not be nil")
		return nil
	}
	if len(job.Payload) == 0 {
		hlog.Fix("ticktock: the len of ticktock job.payload should not be 0")
		return nil
	}
	innerCtx := hlog.NewTraceCtx(ctx)
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(innerCtx, "ticktock.periodicJobHandle",
			hlog.F(zap.String("job.type", job.Type)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err), zap.String("task.payload", string(job.Payload))}
				}
				return []zap.Field{}
			},
		)()
	}
	periodicJobPayload := hjson.MustFromBytes[PeriodicJobPayload](job.Payload)
	if periodicJobPayload == nil {
		hlog.TraceFix("ticktock: the PeriodicJobPayload should not be nil",
			ctx, fmt.Errorf("ticktock: the PeriodicJobPayload should not be nil"))
		return nil
	}

	if job.Type == MqTaskType {
		nxtTime, err := periodicJobPayload.Expression.Next(periodicJobPayload.Current)
		if err != nil {
			hlog.TraceFix("ticktock.onPeriodicJobHandlerFunc", ctx, err)
			return nil
		}
		mqPayload := hjson.MustFromBytes[MqJobPayload](periodicJobPayload.Payload)
		if mqPayload == nil {
			hlog.TraceFix("ticktock.onPeriodicJobHandlerFunc", ctx, fmt.Errorf("mqPayload is nil"))
			return nil
		}
		LocalSchedule(nxtTime, func() {
			_ = onMqJobHandlerFunc(ctx, &Job{
				Type:    MqTaskType,
				Payload: periodicJobPayload.Payload,
			})
		})
	} else {
		nxtJob := PeriodicJobFromPayload(periodicJobPayload)
		err = w.localPostman.Send(ctx, nxtJob)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Worker) doStartup() (context.Context, error) {
	// todo will change for service dep
	tickTockRunning := hcfg.GetBool("helix.ticktock.running", true)
	if !tickTockRunning {
		return nil, nil
	}

	var err error
	go func() {
		if err = w.srv.Run(w.srvMux); err != nil {
			return
		}
	}()
	return w.baseCtx, nil
}

func (w *Worker) doShutdown(_ context.Context) {
	// todo will change for service dep
	tickTockRunning := hcfg.GetBool("helix.ticktock.running", true)
	if !tickTockRunning {
		return
	}
	if w.stopFunc != nil {
		w.stopFunc()
	}
	w.srv.Shutdown()
}

func (w *Worker) HandleError(ctx context.Context, task *asynq.Task, err error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				hlog.Err("ticktock.HandleErr: panic", hlog.TraceInfo(ctx), zap.Any("r", r))
			}
		}()
		hlog.TraceErr("ticktock.HandleErr", ctx, err, zap.String("task_type", task.Type()))
	}()
}

func (w *Worker) cKey(key string) string {
	return fmt.Sprintf("helix.ticktock.%s.%s", strings.ToLower(w.code), key)
}
