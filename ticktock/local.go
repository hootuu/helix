package ticktock

import "time"

func LocalSchedule(target time.Time, job func()) {
	duration := time.Until(target)
	if duration <= 0 {
		go job()
		return
	}

	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		job()
	}()
}
