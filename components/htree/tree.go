package htree

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/honce"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type Tree struct {
	Code    string
	factory *IdFactory
}

func newTree(code string, version uint, cfg []uint) (*Tree, error) {
	f, err := NewIdFactoryWithVersion(version, cfg)
	if err != nil {
		return nil, err
	}
	return &Tree{
		Code:    code,
		factory: f,
	}, nil
}

func (t *Tree) Helix() helix.Helix {
	return helix.BuildHelix(t.Code, func() (context.Context, error) {
		err := t.doInit()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	})
}

func (t *Tree) Factory() *IdFactory {
	return t.factory
}

func (t *Tree) Root() ID {
	return t.factory.Root()
}

func (t *Tree) Next(id ID, call func(id ID) error) error {
	treeM, err := hdb.Get[TreeM](zplt.HelixPgDB().PG().Table(t.tableName()),
		"id = ?", id)
	if err != nil {
		hlog.Err("helix.tree.Next: Get", zap.Error(err))
		return err
	}
	if treeM == nil {
		hlog.Err("helix.tree.Next: no such id", zap.Int64("id", id))
		return fmt.Errorf("no such id: %d", id)
	}
	nxtId, err := t.factory.Next(id, treeM.Sequence+1)
	if err != nil {
		hlog.Err("helix.tree.Next: f.Next", zap.Error(err))
		return err
	}
	err = hdb.Tx(zplt.HelixPgDB().PG().Table(t.tableName()), func(tx *gorm.DB) error {
		err := hdb.Update[TreeM](tx,
			map[string]interface{}{
				"sequence": gorm.Expr("sequence + 1"),
				"version":  gorm.Expr("version + 1"),
			},
			"id = ?", id)
		if err != nil {
			hlog.Err("helix.tree.Next: Update", zap.Error(err))
			return err
		}
		newTreeM := &TreeM{
			ID:       nxtId,
			Sequence: 0,
			Version:  0,
		}
		err = hdb.Create[TreeM](tx, newTreeM)
		if err != nil {
			hlog.Err("helix.tree.Next: Create", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = call(nxtId)
	if err != nil {
		tx := zplt.HelixPgDB().PG().Table(t.tableName()).Unscoped().Delete(&TreeM{}, "id = ?", nxtId)
		if tx.Error != nil {
			hlog.Err("[ignore]helix.tree.Next: Delete", zap.Error(tx.Error), zap.Int64("id", nxtId))
		}
		return err
	}

	return nil
}

func (t *Tree) NextID(parent ID) (ID, error) {
	var id ID
	err := t.Next(parent, func(newID ID) error {
		id = newID
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (t *Tree) doInit() error {
	err := hdb.AutoMigrateWithTable(zplt.HelixPgDB().PG(), hdb.NewTable(t.tableName(), &TreeM{}))
	if err != nil {
		hlog.Err("helix.tree.doInit: AutoMigrateWithTable", zap.Error(err))
		return err
	}
	err = honce.Do(fmt.Sprintf("helix.tree.%s.%d.root.init",
		strings.ToLower(t.Code), t.factory.version),
		func() error {
			treeM, err := hdb.Get[TreeM](zplt.HelixPgDB().PG().Table(t.tableName()),
				"id = ?", t.Root())
			if err != nil {
				hlog.Err("helix.tree.doInit: Get", zap.Error(err))
				return err
			}
			if treeM == nil {
				rootTreeM := &TreeM{
					ID:       t.Root(),
					Sequence: 0,
					Version:  0,
				}
				err = hdb.Create[TreeM](zplt.HelixPgDB().PG().Table(t.tableName()), rootTreeM)
				if err != nil {
					hlog.Err("helix.tree.doInit: Create", zap.Error(err))
					return err
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

func (t *Tree) tableName() string {
	return fmt.Sprintf("helix_tree_%s", strings.ToLower(t.Code))
}
