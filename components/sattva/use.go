package sattva

import (
	"context"
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/components/sattva/attribute"
	"github.com/hootuu/helix/components/sattva/channel"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func NewSattva(code string) (*Sattva, error) {
	if err := helix.CheckCode(code); err != nil {
		return nil, err
	}
	return newSattva(code), nil
}

var gUniIdGenerator hnid.Generator

func init() {
	helix.Use(helix.BuildHelix("helix_sattva_uni", func() (context.Context, error) {
		RegisterBuilder(Password, channel.NewPwdBuilder())
		var err error
		err = zplt.HelixPgDB().PG().AutoMigrate(
			&ChannelM{},
			&IdentificationM{},
			&channel.IdChannelM{},
			&attribute.SimpleM{},
			&attribute.ComplexM{},
		)
		if err != nil {
			hlog.Err("sattva.AutoMigrate", zap.Error(err))
			return nil, err
		}
		gUniIdGenerator, err = hnid.NewGenerator(
			"helix_sattva_uni_id_generator",
			hnid.NewOptions(1, 8).
				SetTimestamp(hnid.Hour, false).
				SetAutoInc(7, 1, 9999999, 20000),
		)
		if err != nil {
			hlog.Err("sattva.gUniIdGenerator.init", zap.Error(err))
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
