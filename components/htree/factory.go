package htree

import (
	"errors"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"math"
)

const (
	maxLength     = 18
	defaultIdRoot = uint(8)
)

type IdFactory struct {
	version uint
	root    ID
	max     ID
	cfg     []uint
	length  int
}

func NewIdFactory(cfg []uint) (*IdFactory, error) {
	return NewIdFactoryWithVersion(defaultIdRoot, cfg)
}

func NewIdFactoryWithVersion(version uint, cfg []uint) (*IdFactory, error) {
	if version == 0 || version > 9 {
		return nil, fmt.Errorf("invalid version, must be [1, 9], but %d", version)
	}
	if len(cfg) == 0 {
		return nil, errors.New("require valid level config")
	}
	sum := 0
	for _, item := range cfg {
		if item == 0 {
			return nil, errors.New("the data length of a single level is equal to 0")
		}
		sum += int(item)
	}
	if sum > maxLength {
		return nil, fmt.Errorf("the total length exceeds the maximum length by %d", maxLength)
	}
	pureRoot := int64(math.Pow10(sum))
	return &IdFactory{
		version: version,
		root:    int64(version) * pureRoot,
		max:     int64(version)*pureRoot + (pureRoot - 1),
		cfg:     cfg,
		length:  sum,
	}, nil
}

func (f *IdFactory) Root() ID {
	return f.root
}

func (f *IdFactory) IsRoot(id ID) bool {
	return id == f.root
}

func (f *IdFactory) IdDeep() int {
	return len(f.cfg)
}

func (f *IdFactory) IdLength() int {
	return f.length
}

func (f *IdFactory) Deep(id ID) int {
	if err := f.validate(id); err != nil {
		return 0
	}
	pure := id - f.root
	left := f.length
	lst := f.root
	for i, item := range f.cfg {
		now := int64(math.Pow10(left - int(item)))
		cur := (pure / now) * now
		if cur == 0 {
			return i
		}
		pure = pure - cur
		left = left - int(item)
		lst = lst + cur
	}
	return f.IdDeep()
}

func (f *IdFactory) Next(id ID, seq int64) (ID, error) {
	minId, maxId, base, err := f.DirectChildren(id)
	if err != nil {
		return 0, err
	}
	nxt := id + seq*base
	if nxt < minId || nxt > maxId {
		hlog.Err("helix.tree.f.Next: invalid seq",
			zap.Int64("seq", seq),
			zap.Int64("min", minId),
			zap.Int64("max", maxId),
			zap.Int64("base", base),
			zap.Int64("nxt", nxt))
		return 0, fmt.Errorf("invalid seq: %d", seq)
	}
	return nxt, nil
}

func (f *IdFactory) Path(id ID) ([]ID, error) {
	if err := f.validate(id); err != nil {
		return nil, err
	}
	pure := id - f.root
	var path []ID
	left := f.length
	lst := f.root
	for _, item := range f.cfg {
		now := int64(math.Pow10(left - int(item)))
		cur := (pure / now) * now
		curID := lst + cur
		path = append(path, curID)
		pure = pure - cur
		left = left - int(item)
		lst = curID
	}
	return path, nil
}

func (f *IdFactory) Children(id ID) (ID, ID, error) {
	if err := f.validate(id); err != nil {
		return 0, 0, err
	}
	cfgLen := len(f.cfg)
	purePath := f.calcPurePath(id)
	if len(purePath) != cfgLen {
		return 0, 0, fmt.Errorf("invalid id: %d, path error: %v [cfg.len: %d]", id, purePath, cfgLen)
	}
	found := false
	idx := 0
	for i := 0; i < cfgLen; i++ {
		if purePath[i] == 0 {
			found = true
			idx = i
			break
		}
	}
	if !found {
		return 0, 0, fmt.Errorf("invalid id: %d, no children", id)
	}
	root := f.root
	for i := 0; i < idx; i++ {
		root += purePath[i]
	}
	powBits := 0
	for i := cfgLen - 1; i >= idx; i-- {
		powBits += int(f.cfg[i])
	}
	maxPow := int64(math.Pow10(powBits))
	minId := root + 1
	maxId := root + maxPow - 1
	return minId, maxId, nil
}

func (f *IdFactory) DirectChildren(id ID) (ID, ID, ID, error) {
	if err := f.validate(id); err != nil {
		return 0, 0, 0, err
	}
	cfgLen := len(f.cfg)
	purePath := f.calcPurePath(id)
	if len(purePath) != cfgLen {
		return 0, 0, 0, fmt.Errorf("invalid id: %d, path error: %v [cfg.len: %d]", id, purePath, cfgLen)
	}
	found := false
	idx := 0
	for i := 0; i < cfgLen; i++ {
		if purePath[i] == 0 {
			found = true
			idx = i
			break
		}
	}
	if !found {
		return 0, 0, 0, fmt.Errorf("invalid id: %d, no children", id)
	}
	root := f.root
	for i := 0; i <= idx; i++ {
		root += purePath[i]
	}
	minBits := 0
	if idx < cfgLen-1 {
		minBits = int(f.cfg[idx+1])
	}
	maxBits := int(f.cfg[idx])
	leftBits := 0
	for i := cfgLen - 1; i > idx+1; i-- {
		leftBits += int(f.cfg[i])
	}
	leftPow := int64(math.Pow10(leftBits))
	minPow := int64(math.Pow10(minBits))
	maxPow := int64(math.Pow10(maxBits))
	basePow := minPow * leftPow
	minId := root + basePow
	maxId := root + (maxPow-1)*basePow
	return minId, maxId, basePow, nil
}

func (f *IdFactory) validate(id ID) error {
	if id < f.root || id > f.max {
		return fmt.Errorf("invalid id: %d", id)
	}
	return nil
}

func (f *IdFactory) calcPurePath(id ID) []int64 {
	pure := id - f.root
	var path []int64
	left := f.length
	lst := f.root
	for _, item := range f.cfg {
		now := int64(math.Pow10(left - int(item)))
		cur := (pure / now) * now
		path = append(path, cur)
		pure = pure - cur
		left = left - int(item)
		lst = lst + cur
	}
	return path
}
