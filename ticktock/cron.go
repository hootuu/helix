package ticktock

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Expression string

var gCronParser = cron.NewParser(
	cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

func (e Expression) ToString() string {
	return string(e)
}

func (e Expression) Validate() error {
	_, err := gCronParser.Parse(string(e))
	if err != nil {
		return err
	}
	return nil
}

func (e Expression) Next(t time.Time) (time.Time, error) {
	schedule, err := gCronParser.Parse(string(e))
	if err != nil {
		return time.Time{}, err
	}
	return schedule.Next(t), nil
}
