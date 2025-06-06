package slice

import (
	"github.com/hootuu/hyle/data/hcast"
	"time"
)

type TimestampType uint8

const (
	YearTT   TimestampType = 1
	MonthTT  TimestampType = 2
	DayTT    TimestampType = 3
	HourTT   TimestampType = 4
	MinuteTT TimestampType = 5
	SecondTT TimestampType = 6
)

const (
	MillPerSecond = int64(1000)
	MillPerMinute = 60 * MillPerSecond
	MillPerHour   = 60 * MillPerMinute
	MillPerDay    = 24 * MillPerHour
	MillPerMonth  = 30 * MillPerDay
	MillPerYear   = 365 * MillPerDay
	MillZero      = int64(1580608922000)
)

type Timestamp struct {
	tt            TimestampType
	useDateFormat bool
}

func NewTimestamp(tt TimestampType, useDateFormat bool) *Timestamp {
	return &Timestamp{
		tt:            tt,
		useDateFormat: useDateFormat,
	}
}

func (ts *Timestamp) Build() (uint64, uint8) {
	if ts.useDateFormat {
		return ts.buildDf()
	}
	return ts.buildUnix()
}

func (ts *Timestamp) buildDf() (uint64, uint8) {
	var layout string
	switch ts.tt {
	case YearTT:
		layout = "2006"
	case MonthTT:
		layout = "200601"
	case DayTT:
		layout = "20060102"
	case HourTT:
		layout = "2006010215"
	case MinuteTT:
		layout = "200601021504"
	case SecondTT:
		layout = "20060102150405"
	default:
		return 0, 0
	}

	str := time.Now().Format(layout)
	long := hcast.ToUint64(str)
	return long, uint8(len(str))
}

func (ts *Timestamp) buildUnix() (uint64, uint8) {
	ms := time.Now().UnixMilli() - MillZero
	numb := uint64(0)
	switch ts.tt {
	case YearTT:
		numb = uint64(ms / MillPerYear)
	case MonthTT:
		numb = uint64(ms / MillPerMonth)
	case DayTT:
		numb = uint64(ms / MillPerDay)
	case HourTT:
		numb = uint64(ms / MillPerHour)
	case MinuteTT:
		numb = uint64(ms / MillPerMinute)
	case SecondTT:
		numb = uint64(ms / MillPerSecond)
	default:
		return 0, 0
	}
	strNumb := hcast.ToString(numb)
	return numb, uint8(len(strNumb))
}
