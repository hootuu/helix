package attribute

import (
	"bytes"
	"errors"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func doSet(id string, attr string, value interface{}) error {
	val := &Value{data: value}
	if !val.IsNil() {
		if val.IsSimple() {
			return doSetSimple(id, attr, value)
		}

		if val.IsComplex() {
			return doSetComplex(id, attr, value)
		}
	}

	return errors.New("the value must be simple[bool/int/string...] or complex[map[string]any]")
}

func doSetSimple(id string, attr string, value interface{}) error {
	strVal := hcast.ToString(value)
	if len(strVal) > 500 {
		return errors.New("the length of attr.value must <= 500")
	}
	simpleM, err := hpg.Get[SimpleM](zplt.HelixPgDB().PG(),
		"identification = ? AND attr = ?",
		id, attr)
	if err != nil {
		hlog.Err("helix.sattva.attr.doSetSimple: hpg.Get", zap.Error(err))
		return err
	}
	if simpleM == nil {
		simpleM = &SimpleM{
			BasicM: BasicM{
				Identification: id,
				Attr:           attr,
			},
			Value: strVal,
		}
		err = hpg.Create[SimpleM](zplt.HelixPgDB().PG(), simpleM)
		if err != nil {
			hlog.Err("helix.sattva.attr.doSetSimple: hpg.Create", zap.Error(err))
			return err
		}
		return nil
	}
	if simpleM.Value == strVal {
		return nil
	}
	mut := map[string]interface{}{
		"value": strVal,
	}
	err = hpg.Update[SimpleM](zplt.HelixPgDB().PG(), mut,
		"identification = ? AND attr = ?",
		id, attr)
	if err != nil {
		hlog.Err("helix.sattva.attr.doSetSimple: hpg.Update", zap.Error(err))
		return err
	}
	return nil
}

func doSetComplex(id string, attr string, value interface{}) error {
	mapVal, _ := value.(map[string]interface{})
	complexM, err := hpg.Get[ComplexM](zplt.HelixPgDB().PG(),
		"identification = ? AND attr = ?",
		id, attr)
	if err != nil {
		hlog.Err("helix.sattva.attr.doSetComplex: hpg.Get", zap.Error(err))
		return err
	}
	mapJsonBytes, err := hjson.ToBytes(mapVal)
	if err != nil {
		hlog.Err("helix.sattva.attr.doSetComplex: hjson.ToBytes", zap.Error(err))
		return err
	}
	if complexM == nil {
		complexM = &ComplexM{
			BasicM: BasicM{
				Identification: id,
				Attr:           attr,
			},
			Value: mapJsonBytes,
		}
		err = hpg.Create[ComplexM](zplt.HelixPgDB().PG(), complexM)
		if err != nil {
			hlog.Err("helix.sattva.attr.doSetComplex: hpg.Create", zap.Error(err))
			return err
		}
		return nil
	}
	if bytes.Equal(mapJsonBytes, hjson.MustToBytes(complexM.Value)) {
		return nil
	}
	mut := map[string]interface{}{
		"value": mapJsonBytes,
	}
	err = hpg.Update[ComplexM](zplt.HelixPgDB().PG(), mut,
		"identification = ? AND attr = ?",
		id, attr)
	if err != nil {
		hlog.Err("helix.sattva.attr.doSetComplex: hpg.Update", zap.Error(err))
		return err
	}
	return nil
}

func doGet(id string, withSimple bool, withComplex bool, attr ...string) (dict.Dict, error) {
	attrDict := dict.NewDict()
	if withSimple {
		simpleArr, err := hpg.Find[SimpleM](func() *gorm.DB {
			return zplt.HelixPgDB().PG().Model(&SimpleM{}).
				Where("identification = ? AND attr in ?", id, attr)
		})
		if err != nil {
			hlog.Err("helix.sattva.attr.doGet: find.simple", zap.Error(err))
			return nil, err
		}
		if len(simpleArr) > 0 {
			for _, m := range simpleArr {
				attrDict.Set(m.Attr, m.Value)
			}
		}
	}
	if withComplex {
		complexArr, err := hpg.Find[ComplexM](func() *gorm.DB {
			return zplt.HelixPgDB().PG().Model(&ComplexM{}).
				Where("identification = ? AND attr in ?", id, attr)
		})
		if err != nil {
			hlog.Err("helix.sattva.attr.doGet: find.complex", zap.Error(err))
			return nil, err
		}
		if len(complexArr) > 0 {
			for _, m := range complexArr {
				ptrMap, err := hjson.FromBytes[map[string]interface{}](m.Value)
				if err != nil {
					hlog.Err("helix.sattva.attr.doGet: complex to json obj", zap.Error(err))
					continue
				}
				attrDict.Set(m.Attr, *ptrMap)
			}
		}
	}
	return attrDict, nil
}
