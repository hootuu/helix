package horg

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/gorm"
)

const (
	aliasMaxLen = 50
	aliasCutLen = 15
)

type CreateParas struct {
	Biz   collar.Collar `json:"biz"`
	Alias string        `json:"alias"`
	Name  string        `json:"name"`
	Meta  dict.Dict     `json:"meta"`
}

func (p CreateParas) Validate() error {
	if p.Biz == "" {
		return errors.New("biz is required")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (p CreateParas) GetAlias() string {
	if p.Alias == "" {
		if len(p.Name) < aliasMaxLen {
			return p.Name
		} else {
			return p.Name[:aliasCutLen] + "..."
		}
	}
	return p.Alias
}

func Create(ctx context.Context, paras CreateParas, call ...func(ctx context.Context, orgM *OrgM) error) (ID, error) {
	if err := paras.Validate(); err != nil {
		return 0, err
	}
	tx := zplt.HelixPgCtx(ctx)
	id, err := gOrgIdTree.NextID(gOrgIdTree.Root())
	if err != nil {
		return 0, err
	}
	orgM := &OrgM{
		Biz:       paras.Biz.ID(),
		Sovereign: true,
		ID:        id,
		Parent:    gOrgIdTree.Root(),
		Alias:     paras.GetAlias(),
		Name:      paras.Name,
		Meta:      hjson.MustToBytes(paras.Meta),
	}
	err = hpg.Create[OrgM](tx, orgM)
	if err != nil {
		return 0, err
	}
	if len(call) > 0 {
		err = call[0](ctx, orgM)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

type AddParas struct {
	Parent    ID        `json:"parent"`
	Sovereign bool      `json:"sovereign"`
	Alias     string    `json:"alias"`
	Name      string    `json:"name"`
	Meta      dict.Dict `json:"meta"`
}

func (p AddParas) Validate() error {
	if p.Parent == 0 {
		return errors.New("parent is required")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (p AddParas) GetAlias() string {
	if p.Alias == "" {
		if len(p.Name) < aliasMaxLen {
			return p.Name
		} else {
			return p.Name[:aliasCutLen] + "..."
		}
	}
	return p.Alias
}

func Add(ctx context.Context, paras AddParas, call ...func(ctx context.Context, orgM *OrgM) error) (ID, error) {
	if err := paras.Validate(); err != nil {
		return 0, err
	}
	tx := zplt.HelixPgCtx(ctx)
	parentM, err := hpg.Get[OrgM](tx, "id = ?", paras.Parent)
	if err != nil {
		return 0, err
	}
	if parentM == nil {
		return 0, errors.New("parent not found")
	}
	id, err := gOrgIdTree.NextID(parentM.ID)
	if err != nil {
		return 0, err
	}
	orgM := &OrgM{
		Biz:       parentM.Biz,
		Sovereign: paras.Sovereign,
		ID:        id,
		Parent:    parentM.ID,
		Alias:     paras.GetAlias(),
		Name:      paras.Name,
		Meta:      hjson.MustToBytes(paras.Meta),
	}
	err = hpg.Create[OrgM](tx, orgM)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func MustExist(ctx context.Context, id ID) error {
	return mustExist(ctx, "id = ?", id)
}

func MustExistWithSovereign(ctx context.Context, id ID, sov bool) error {
	return mustExist(ctx, "id = ? AND sovereign = ?", id, sov)
}

func mustExist(ctx context.Context, query any, cond ...any) error {
	b, err := hpg.Exist[OrgM](zplt.HelixPgCtx(ctx), query, cond...)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("org not exist: %s", hjson.MustToString(cond))
	}
	return nil
}

func Get(ctx context.Context, parent ID, deep int) ([]*Organization, error) {
	if deep < 1 || deep > gOrgIdTree.Factory().IdDeep() {
		return nil, fmt.Errorf("invalid deep: %d", deep)
	}
	if parent == Root {
		parent = gOrgIdTree.Root()
	}
	minID, maxID, base, err := gOrgIdTree.Factory().DirectChildren(parent)
	if err != nil {
		return nil, err
	}
	var arr []*Organization
	arr, err = loadChildren(ctx, minID, maxID, base)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return []*Organization{}, nil
	}
	newDeep := deep - 1
	if newDeep <= 0 {
		return arr, nil
	}
	for _, categ := range arr {
		categ.Children, err = Get(ctx, categ.ID, newDeep)
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func loadChildren(ctx context.Context, minID htree.ID, maxID htree.ID, base htree.ID) ([]*Organization, error) {
	arrM, err := hpg.Find[OrgM](func() *gorm.DB {
		return zplt.HelixPgCtx(ctx).Where("id % ? = 0 AND id >= ? AND id <= ?", base, minID, maxID)
	})
	if err != nil {
		return []*Organization{}, err
	}
	if len(arrM) == 0 {
		return []*Organization{}, nil
	}
	var arr []*Organization
	for _, item := range arrM {
		arr = append(arr, item.ToOrganization())
	}
	return arr, nil
}
