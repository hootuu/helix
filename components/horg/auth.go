package horg

import (
	"context"
	"errors"
	"github.com/hootuu/helix/components/hseq"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func AddAuth(ctx context.Context, orgID ID, name string, action []string) (authID AuthID, err error) {
	defer hlog.ElapseWithCtx(ctx, "horg.AddAuth", func() []zap.Field {
		return []zap.Field{zap.Int64("orgID", orgID)}
	}, func() []zap.Field {
		if err != nil {
			return []zap.Field{zap.Error(err)}
		}
		return []zap.Field{}
	})()

	if orgID == 0 {
		return 0, errors.New("require orgID")
	}
	if name == "" {
		return 0, errors.New("require name")
	}
	if len(action) == 0 {
		return 0, errors.New("require action")
	}
	if err := MustExistWithSovereign(ctx, orgID, true); err != nil {
		return 0, err
	}
	authID, err = hseq.Next(ctx, CollarAuth(orgID))
	if err != nil {
		return 0, err
	}
	tx := zplt.HelixPgCtx(ctx)
	m := &OrgAuthorityM{
		InnerID: authID,
		Org:     orgID,
		Name:    name,
		Action:  action,
	}
	err = hdb.Create[OrgAuthorityM](tx, m)
	if err != nil {
		return 0, err
	}
	return authID, nil
}

func SetAuth(ctx context.Context, orgID ID, authID AuthID, name string, action []string) error {
	if orgID == 0 {
		return errors.New("require orgID")
	}
	if authID == 0 {
		return errors.New("require authID")
	}
	tx := zplt.HelixPgCtx(ctx)
	authM, err := hdb.Get[OrgAuthorityM](tx, "org = ? AND inner_id = ?", orgID, authID)
	if err != nil {
		return err
	}
	if authM == nil {
		return errors.New("auth not found")
	}
	mut := make(map[string]any)
	if name != "" && authM.Name != name {
		mut["name"] = name
	}
	if len(action) > 0 {
		mut["action"] = action
	}
	err = hdb.Update[OrgAuthorityM](tx, mut, "org = ? AND inner_id = ?", orgID, authID)
	if err != nil {
		return err
	}
	return nil
}
