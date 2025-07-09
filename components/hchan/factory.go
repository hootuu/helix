package hchan

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/components/zplt/zcanal"
	"github.com/hootuu/helix/components/zplt/zmeili"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hcanal"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type Factory struct {
	Code string
	tree *htree.Tree
	db   *hdb.Database
}

func newFactory(
	code string,
	flag uint,
	cfg []uint,
) (*Factory, error) {
	c := &Factory{
		Code: code,
		db:   zplt.HelixDB(),
		tree: nil,
	}
	var err error
	c.tree, err = htree.NewTree(fmt.Sprintf("hchan_%s", code), flag, cfg)
	if err != nil {
		hlog.Err("hchan.newFactory: NewTree", zap.Error(err))
		return nil, err
	}
	return c, nil
}

func (f *Factory) Root() ID {
	return f.tree.Root()
}

func (f *Factory) Add(parent ID, name string, icon string, seq int) (ID, error) {
	if name == "" {
		return 0, errors.New("require name")
	}
	if parent == Root {
		parent = f.tree.Root()
	}
	b, err := hdb.Exist[ChanM](f.table(), "parent = ? AND name = ?", parent, name)
	if err != nil {
		return -1, err
	}
	if b {
		return -1, fmt.Errorf("exists: parent=%d,name=%s", parent, name)
	}
	var newID htree.ID
	err = f.tree.Next(parent, func(id htree.ID) error {
		newID = id
		return nil
	})
	if err != nil {
		return -1, err
	}
	chanM := &ChanM{
		ID:        newID,
		Parent:    parent,
		Name:      name,
		Icon:      icon,
		Seq:       seq,
		Available: true,
	}
	err = hdb.Create[ChanM](f.table(), chanM)
	if err != nil {
		return -1, err
	}
	return chanM.ID, nil
}

func (f *Factory) Mut(id ID, name string, icon string, seq int) error {
	if name == "" {
		return errors.New("require name")
	}
	dbM, err := hdb.MustGet[ChanM](f.table(), "id = ?", id)
	if err != nil {
		return err
	}
	mut := make(map[string]any)
	if dbM.Name != name {
		mut["name"] = name
	}
	if dbM.Icon != icon {
		mut["icon"] = icon
	}
	if dbM.Seq != seq {
		mut["seq"] = seq
	}
	if len(mut) == 0 {
		return nil
	}
	b, err := hdb.Exist[ChanM](f.table(), "parent = ? AND name = ? AND id <> ?", dbM.Parent, name, id)
	if err != nil {
		return err
	}
	if b {
		return fmt.Errorf("exists: parent=%d,name=%s", dbM.Parent, name)
	}

	err = hdb.Update[ChanM](f.table(), mut, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) SetAvailable(id ID, available bool) error {
	dbM, err := hdb.MustGet[ChanM](f.table(), "id = ?", id)
	if err != nil {
		return err
	}
	if dbM.Available == available {
		return nil
	}
	mut := map[string]any{
		"available": available,
	}
	err = hdb.Update[ChanM](f.table(), mut, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) SetSeq(id ID, seq int) error {
	dbM, err := hdb.MustGet[ChanM](f.table(), "id = ?", id)
	if err != nil {
		return err
	}
	if dbM.Seq == seq {
		return nil
	}
	mut := map[string]any{
		"seq": seq,
	}
	err = hdb.Update[ChanM](f.table(), mut, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (f *Factory) GetChildren(parent ID, deep int) ([]*Channel, error) {
	if deep < 1 || deep > f.tree.Factory().IdDeep() {
		return nil, fmt.Errorf("invalid deep: %d", deep)
	}
	if parent == Root {
		parent = f.tree.Root()
	}
	minID, maxID, base, err := f.tree.Factory().DirectChildren(parent)
	if err != nil {
		return nil, err
	}
	var arr []*Channel
	arr, err = f.loadChildren(minID, maxID, base)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return []*Channel{}, nil
	}
	newDeep := deep - 1
	if newDeep <= 0 {
		return arr, nil
	}
	for _, ch := range arr {
		ch.Children, err = f.GetChildren(ch.ID, newDeep)
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func (f *Factory) Filter(filter string, sort []string, page *pagination.Page) (*pagination.Pagination[any], error) {
	return hmeili.Filter(zmeili.HelixMeili(), f.tableName(), filter, sort, page)
}

func (f *Factory) loadChildren(minID htree.ID, maxID htree.ID, base htree.ID) ([]*Channel, error) {
	arrM, err := hdb.Find[ChanM](func() *gorm.DB {
		return f.table().Where("id % ? = 0 AND id >= ? AND id <= ?", base, minID, maxID)
	})
	if err != nil {
		return []*Channel{}, err
	}
	if len(arrM) == 0 {
		return []*Channel{}, nil
	}
	var arr []*Channel
	for _, item := range arrM {
		arr = append(arr, item.To())
	}
	return arr, nil
}

func (f *Factory) Helix() helix.Helix {
	return helix.BuildHelix(f.Code, func() (context.Context, error) {
		return f.doStartup()
	}, func(ctx context.Context) {

	})
}

func (f *Factory) table() *gorm.DB {
	return f.db.DB().Table(f.tableName())
}

func (f *Factory) tableName() string {
	return fmt.Sprintf("helix_hchan_%s", strings.ToLower(f.Code))
}

func (f *Factory) doStartup() (context.Context, error) {
	err := hdb.AutoMigrateWithTable(f.db.DB(), hdb.NewTable(f.tableName(), &ChanM{}))
	if err != nil {
		hlog.Err("helix.hchan.doStartup: AutoMigrateWithTable", zap.Error(err))
		return nil, err
	}
	meiliPtr := zmeili.HelixMeili()
	indexer := newChannelIndexer(f)
	err = hmeili.InitIndexer(meiliPtr, indexer)
	if err != nil {
		hlog.Err("helix.hchan.doStartup: init indexer failed",
			zap.String("idx", indexer.GetName()), zap.Error(err))
		return nil, err
	}
	zcanal.HelixCanal().RegisterAlterHandler(
		hcanal.NewIndexHandler(f.tableName(), indexer, meiliPtr),
	)
	return nil, nil
}
