package hwt

import (
	"context"
	"errors"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/helix/storage/hrds"
	"github.com/hootuu/hyle/crypto/hed25519"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hio"
	"github.com/hootuu/hyle/hlog"
	"github.com/mr-tron/base58"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Hwt struct {
	Code       string
	refreshExp time.Duration
	tokenExp   time.Duration
}

func newHwt(code string) (*Hwt, error) {
	hwt := &Hwt{
		Code:       code,
		refreshExp: hcfg.GetDuration("hwt."+code+".refresh.expiration", 7*24*time.Hour),
		tokenExp:   hcfg.GetDuration("hwt."+code+".token.expiration", 24*time.Hour),
	}
	helix.Use(hwt.Helix())
	return hwt, nil
}

func (h *Hwt) Helix() helix.Helix {
	return helix.BuildHelix(
		"helix_hwt_"+h.Code,
		func() (context.Context, error) {
			err := hpg.AutoMigrateWithTable(h.db().PG(),
				hpg.NewTable(buildRefreshTokenTableName(h.Code), &RefreshTokenM{}))
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
		func(ctx context.Context) {

		},
	)
}

func (h *Hwt) RefreshIssuing(identity string) (*hio.JwtToken, error) {
	current := time.Now()
	var jwtToken hio.JwtToken
	var err error

	jwtToken.Refresh, err = buildToken()
	if err != nil {
		hlog.Err("hwt.RefreshIssuing: buildToken", zap.Error(err))
		return nil, err
	}
	jwtToken.RefreshExpiration = current.Add(h.refreshExp).UnixMilli()

	err = hpg.Tx(
		h.db().PG().Table(buildRefreshTokenTableName(h.Code)),
		func(tx *gorm.DB) error {
			var err error
			//delete old
			err = tx.Unscoped().Where(
				"identity = ? AND code = ?",
				identity, h.Code,
			).Delete(&RefreshTokenM{}).Error
			if err != nil {
				hlog.Err("hwt.RefreshIssuing: Delete Old", zap.Error(err))
				return err
			}

			//create new
			model := &RefreshTokenM{
				RefreshToken: jwtToken.Refresh,
				Identity:     identity,
				Code:         h.Code,
				Expiration:   jwtToken.RefreshExpiration,
			}
			err = tx.Create(model).Error
			if err != nil {
				hlog.Err("hwt.RefreshIssuing: Create New", zap.Error(err))
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	jwtToken.Token, err = h.doTokenIssuing(identity)
	if err != nil {
		return nil, err
	}
	jwtToken.TokenExpiration = current.Add(h.tokenExp).UnixMilli()
	return &jwtToken, nil
}

func (h *Hwt) TokenIssuing(identity string, refreshToken string) (*hio.JwtToken, error) {
	current := time.Now()
	var jwtToken hio.JwtToken
	var err error

	refreshM, err := hpg.Get[RefreshTokenM](
		h.db().PG().Table(buildRefreshTokenTableName(h.Code)),
		"refresh_token = ? AND identity = ? AND code = ?",
		refreshToken, identity, h.Code,
	)
	if err != nil {
		hlog.Err("hwt.TokenIssuing: Create New", zap.Error(err))
		return nil, err
	}
	if refreshM == nil {
		return nil, errors.New("refresh token not exists or have expired")
	}
	if current.After(time.UnixMilli(refreshM.Expiration)) {
		return nil, errors.New("refresh token have expired")
	}
	jwtToken.Refresh = refreshToken
	jwtToken.RefreshExpiration = refreshM.Expiration
	jwtToken.Token, err = h.doTokenIssuing(identity)
	if err != nil {
		hlog.Err("hwt.TokenIssuing: buildToken", zap.Error(err))
		return nil, err
	}
	jwtToken.TokenExpiration = current.Add(h.tokenExp).UnixMilli()
	return &jwtToken, nil
}

func (h *Hwt) doTokenIssuing(identity string) (string, error) {
	token, err := buildToken()
	if err != nil {
		hlog.Err("hwt.doTokenIssuing", zap.Error(err))
		return "", err
	}
	err = h.cache().Set(token, identity, h.tokenExp)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (h *Hwt) db() *hpg.Database {
	return zplt.HelixPgDB()
}

func (h *Hwt) cache() *hrds.Cache {
	return zplt.HelixRdsCache()
}

func buildToken() (string, error) {
	_, pri, err := hed25519.Random()
	if err != nil {
		return "", err
	}
	return base58.Encode(pri), nil
}
