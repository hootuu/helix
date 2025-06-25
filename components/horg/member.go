package horg

import (
	"context"
	"errors"
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func BindMember(
	ctx context.Context,
	orgID ID,
	memberID sattva.Identification,
	authority []AuthID,
) (err error) {
	defer hlog.ElapseWithCtx(ctx, "horg.BindMember", func() []zap.Field {
		return []zap.Field{
			zap.Int64("orgID", orgID),
			zap.String("memberID", memberID),
		}
	}, func() []zap.Field {
		if err != nil {
			return []zap.Field{zap.Error(err)}
		}
		return []zap.Field{}
	})()

	if orgID == 0 {
		return errors.New("require orgID")
	}
	if memberID == "" {
		return errors.New("require memberID")
	}
	if err := MustExist(ctx, orgID); err != nil {
		return err
	}
	tx := zplt.HelixPgCtx(ctx)

	memberM, err := hdb.Get[OrgMemberM](tx, "org = ? AND member = ?", orgID, memberID)
	if err != nil {
		return err
	}
	if memberM == nil {
		memberM = &OrgMemberM{
			Org:       orgID,
			Member:    memberID,
			Authority: authority,
		}
		err = hdb.Create[OrgMemberM](tx, memberM)
		if err != nil {
			return err
		}
		return nil
	}
	if len(authority) > 0 {
		mut := map[string]any{
			"authority": authority,
		}
		err = hdb.Update[OrgMemberM](tx, mut, "org = ? AND member = ?", orgID, memberID)
		if err != nil {
			return err
		}
	}
	return nil
}
