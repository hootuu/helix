package sattva

import (
	"context"
	"github.com/hootuu/helix/components/sattva/attribute"
	"github.com/hootuu/helix/components/sattva/channel"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
)

type Sattva struct {
	Code string
	db   *hdb.Database
}

func newSattva(code string) *Sattva {
	s := &Sattva{
		Code: code,
	}
	helix.Use(s.Helix())
	return s
}

func (s *Sattva) WithDatabase(db *hdb.Database) *Sattva {
	s.db = db
	return s
}

func (s *Sattva) Helix() helix.Helix {
	return helix.BuildHelix("sattva_"+s.Code, func() (context.Context, error) {
		err := s.doInit()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	})
}

func (s *Sattva) RegisterChannel(chnType channel.Type, chnCode string, cfg channel.Config) (channel.ID, error) {
	return uniChannelRegister(chnType, chnCode, cfg)
}

func (s *Sattva) IdentificationCreate(chn *channel.Channel, info dict.Dict) (Identification, error) {
	return uniIdentificationCreate(chn, info)
}

func (s *Sattva) Identify(chn *channel.Channel) (bool, Identification, error) {
	return uniIdentify(chn)
}

func (s *Sattva) SetInfo(id Identification, info dict.Dict) error {
	return uniSetInfo(id, info)
}

func (s *Sattva) GetInfo(id Identification) (dict.Dict, error) {
	return uniGetInfo(id)
}

func (s *Sattva) GetAttribute(id Identification, attr ...string) (dict.Dict, error) {
	return attribute.Get(id, attr...)
}

func (s *Sattva) GetAttrSimple(id Identification, attr ...string) (dict.Dict, error) {
	return attribute.GetSimple(id, attr...)
}

func (s *Sattva) GetAttrComplex(id Identification, attr ...string) (dict.Dict, error) {
	return attribute.GetComplex(id, attr...)
}

func (s *Sattva) SetAttribute(id Identification, attr string, value interface{}) error {
	return attribute.Set(id, attr, value)
}

func (s *Sattva) doInit() error {
	if s.db == nil {
		s.db = zplt.HelixPgDB()
	}
	return nil
}
