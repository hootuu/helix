package attribute

import "github.com/hootuu/hyle/data/dict"

func Set(id string, attr string, value interface{}) error {
	return doSet(id, attr, value)
}

func Get(id string, attr ...string) (dict.Dict, error) {
	return doGet(id, true, true, attr...)
}

func GetSimple(id string, attr ...string) (dict.Dict, error) {
	return doGet(id, true, false, attr...)
}

func GetComplex(id string, attr ...string) (dict.Dict, error) {
	return doGet(id, false, true, attr...)
}
