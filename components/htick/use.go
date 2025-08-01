package htick

import (
	"errors"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"go.uber.org/zap"
)

func Schedule(job *Job) error {
	if job == nil {
		return errors.New("job is required")
	}
	id, err := gCron.AddFunc(string(job.Expression), func() {
		gLogger.Info("htick.Schedule.Run",
			zap.String("topic", string(job.Topic)),
			zap.String("cron", string(job.Expression)))
		err := hretry.Must(func() error {
			err := triggerTick(job)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			hlog.Fix("htick.Schedule.Job", zap.Any("job", job), zap.Error(err))
		}
	})
	if err != nil {
		return err
	}
	gLogger.Info("htick.Schedule", zap.Int("EntryID", int(id)))
	return nil
}
