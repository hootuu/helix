package hnid

import (
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/hnid/slice"
	"math"
	"regexp"
	"sync"
)

type Generator interface {
	Next() NID
	NextString() string
	NextUint64() uint64
}

type TimestampType = slice.TimestampType

const (
	Year   = slice.YearTT
	Month  = slice.MonthTT
	Day    = slice.DayTT
	Hour   = slice.HourTT
	Minute = slice.MinuteTT
	Second = slice.SecondTT
)

type Options struct {
	useTimestamp           bool
	timestampType          TimestampType
	timestampUseDateFormat bool
	bizLen                 uint8
	bizRef                 uint64
	autoIncLen             uint8
	autoIncStart           uint64
	autoIncEnd             uint64
	autoIncStep            uint64
}

func (opt *Options) validate() error {
	if opt.useTimestamp {
		switch opt.timestampType {
		case Year, Month, Day, Hour, Minute, Second:
		default:
			return errors.New(fmt.Sprintf("invalid timestamp type: %d", opt.timestampType))
		}
	}
	if opt.bizLen == 0 || opt.bizLen > 6 {
		return errors.New(fmt.Sprintf("invalid biz len: (0, 6], but: %d", opt.bizLen))
	}
	bizRefMax := uint64(math.Pow10(int(opt.bizLen))) - 1
	if opt.bizRef > bizRefMax {
		return errors.New(fmt.Sprintf("invalid biz ref: biz.ref=%d; biz.len=%d", opt.bizRef, opt.bizLen))
	}
	if opt.autoIncLen == 0 || opt.autoIncLen > 16 {
		return errors.New(fmt.Sprintf("invalid biz len: (0, 16], but: %d", opt.autoIncLen))
	}
	autoIncMax := uint64(math.Pow10(int(opt.autoIncLen))) - 1
	if opt.autoIncStep > autoIncMax {
		return errors.New(fmt.Sprintf("invalid auto inc step: %d > %d", opt.autoIncStep, autoIncMax))
	}
	if opt.autoIncStart > opt.autoIncEnd {
		return errors.New(fmt.Sprintf("invalid auto inc config: %d > %d", opt.autoIncStart, opt.autoIncEnd))
	}
	if opt.autoIncEnd > autoIncMax {
		return errors.New(fmt.Sprintf("invalid auto inc end: %d > %d", opt.autoIncEnd, autoIncMax))
	}
	return nil
}

func (opt *Options) SetTimestamp(timestampType TimestampType, useDateFormat bool) *Options {
	opt.useTimestamp = true
	opt.timestampType = timestampType
	opt.timestampUseDateFormat = useDateFormat
	return opt
}

func (opt *Options) SetAutoInc(autoIncLen uint8, start uint64, end uint64, step uint64) *Options {
	opt.autoIncLen = autoIncLen
	opt.autoIncStart = start
	opt.autoIncEnd = end
	opt.autoIncStep = step
	return opt
}

func NewOptions(bizLen uint8, bizRef uint64) *Options {
	return &Options{
		useTimestamp:           false,
		timestampType:          Second,
		timestampUseDateFormat: true,
		bizLen:                 bizLen,
		bizRef:                 bizRef,
		autoIncLen:             9,
		autoIncStart:           1,
		autoIncEnd:             999999999,
		autoIncStep:            10000,
	}
}

var gGeneratorMap = make(map[string]Generator)
var gGeneratorMapMu sync.Mutex

func NewGenerator(code string, opt *Options) (Generator, error) {
	if err := validateCode(code); err != nil {
		return nil, err
	}
	if err := opt.validate(); err != nil {
		return nil, err
	}
	g, err := doNewLocalGenerator(code, opt)
	if err != nil {
		return nil, err
	}
	gGeneratorMapMu.Lock()
	defer gGeneratorMapMu.Unlock()
	gGeneratorMap[code] = g
	return g, nil
}

func Each(call func(code string, g Generator)) {
	for code, g := range gGeneratorMap {
		call(code, g)
	}
}

var gCodeRegexpTpl = `^[A-Za-z0-9_.-]{0,63}$`
var gCodeRegexp = regexp.MustCompile(gCodeRegexpTpl)

func validateCode(code string) error {
	matched := gCodeRegexp.MatchString(code)
	if !matched {
		return errors.New("invalid nid code[" + gCodeRegexpTpl + "]: " + code)
	}
	return nil
}
