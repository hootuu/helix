package htree

import (
	"errors"
	"fmt"
	"math"
)

const (
	maxLength     = 18
	defaultIdRoot = uint(8)
)

type IdFactory struct {
	root   ID
	cfg    []uint
	length uint
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
	sum := uint(0)
	for _, item := range cfg {
		if item == 0 {
			return nil, errors.New("the data length of a single level is equal to 0")
		}
		sum += item
	}
	if sum > maxLength {
		return nil, fmt.Errorf("the total length exceeds the maximum length by %d", maxLength)
	}
	return &IdFactory{root: ID(int64(version) * int64(math.Pow10(int(sum)))), cfg: cfg, length: sum}, nil
}

func (f *IdFactory) Root() ID {
	return f.root
}

func (f *IdFactory) IsRoot(id ID) bool {
	return id == f.root
}

func (f *IdFactory) Deep() int {
	return len(f.cfg)
}

func (f *IdFactory) Path(id ID) ([]ID, error) {
	if err := f.validate(id); err != nil {
		return nil, err
	}
	pure := int64(id - f.root)
	var path []ID
	left := f.length
	lst := f.root
	for _, item := range f.cfg {
		now := int64(math.Pow10(int(left - item)))
		cur := (pure / now) * now
		curID := lst + ID(cur)
		path = append(path, curID)
		pure = pure - cur
		left = left - item
		lst = curID
	}
	return path, nil
}

func (f *IdFactory) DirectChildren(id ID) (ID, ID, error) {
	if err := f.validate(id); err != nil {
		return 0, 0, err
	}
	return 0, 0, nil
}

func (f *IdFactory) validate(id ID) error {
	if id < f.root {
		return fmt.Errorf("invalid id: %d", id)
	}
	return nil
}
