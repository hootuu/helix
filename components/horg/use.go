package horg

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func init() {
	helix.Use(helix.BuildHelix("helix_org", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&OrgM{},
			&OrgAuthorityM{},
			&OrgMemberM{},
		)
		if err != nil {
			return nil, err
		}
		err = initOrgIdTree()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
