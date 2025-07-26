package sattva

import (
	"github.com/hootuu/helix/storage/hrds"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"time"
)

type Summary struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobi     string `json:"mobi"`
}

func FastGet(cache *hrds.Cache, s *Sattva, id Identification) *Summary {
	info, err := hrds.Fast[Summary](
		cache,
		"helix_sattva",
		id,
		6*time.Hour,
		func() (*Summary, error) {
			data, err := s.GetAttrSimple(id, "avatar", "nick_name", "phone")
			if err != nil {
				return nil, err
			}
			return &Summary{
				Nickname: data.Get("nick_name").String(),
				Avatar:   data.Get("avatar").String(),
				Mobi:     data.Get("phone").String(),
			}, nil
		})
	if err != nil {
		hlog.Err("helix.sattva.FastGet",
			zap.String("id", id),
			zap.Error(err),
		)
	}
	if info == nil {
		info = &Summary{
			Nickname: "",
			Avatar:   "",
			Mobi:     "",
		}
	}
	return info
}
