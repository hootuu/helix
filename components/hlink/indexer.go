package hlink

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/components/zplt/zmeili"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const (
	LindIndex      = "helix_link"
	linkIdxVersion = "1.0.0"
)

func Filter(filter string, sort []string, page *pagination.Page) (*pagination.Pagination[any], error) {
	return hmeili.Filter(zmeili.HelixMeili(), LindIndex, filter, sort, page)
}

type linkIndexer struct {
}

func (idx *linkIndexer) GetName() string {
	return LindIndex
}

func (idx *linkIndexer) GetVersion() string {
	return linkIdxVersion
}

func (idx *linkIndexer) Setting(index meilisearch.IndexManager) error {
	filterableAttributes := []string{
		"auto_id",
		"id",
		"major_code",
		"major_id",
		"relation",
		"counterpart_code",
		"counterpart_id",
		"biz",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		hlog.Err("hlink.idx.Setting: Error updating filterable attributes", zap.Error(err))
		return err
	}

	sortableAttributes := []string{
		"auto_id",
		"timestamp",
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		hlog.Err("hlink.idx.Setting: Error updating sortable attributes", zap.Error(err))
		return err
	}
	return nil
}

func (idx *linkIndexer) Load(autoID int64) (hmeili.Document, error) {
	m, err := hdb.MustGet[LinkM](zplt.HelixDB().DB(), "auto_id = ?", autoID)
	if err != nil {
		hlog.Err("hlink.idx.Load: Error loading link", zap.Error(err))
		return nil, err
	}
	doc := hmeili.NewMapDocument(m.ID, m.AutoID, m.UpdatedAt.UnixMilli())
	doc["biz"] = m.Biz
	doc["relation"] = m.Relation
	major, _ := collar.FromID(m.Major)
	majorCode, majorId := major.Parse()
	doc["major_code"] = majorCode
	doc["major_id"] = majorId
	counterpart, _ := collar.FromID(m.Counterpart)
	counterpartCode, counterpartId := counterpart.Parse()
	doc["counterpart_code"] = counterpartCode
	doc["counterpart_id"] = counterpartId
	return doc, nil
}
