package hchan

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const (
	idxVersion = "1.0.0"
)

type channelIndexer struct {
	f *Factory
}

func newChannelIndexer(f *Factory) *channelIndexer {
	return &channelIndexer{f: f}
}

func (idx *channelIndexer) GetName() string {
	return idx.f.tableName()
}

func (idx *channelIndexer) GetVersion() string {
	return idxVersion
}

func (idx *channelIndexer) Setting(index meilisearch.IndexManager) error {
	filterableAttributes := []string{
		"auto_id",
		"id",
		"parent",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		hlog.Err("hchan.idx.Setting: Error updating filterable attributes", zap.Error(err))
		return err
	}

	sortableAttributes := []string{
		"auto_id",
		"timestamp",
		"seq",
		"id",
		"parent",
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		hlog.Err("hchan.idx.Setting: Error updating sortable attributes", zap.Error(err))
		return err
	}
	return nil
}

func (idx *channelIndexer) Load(autoID int64) (hmeili.Document, error) {
	m, err := hdb.MustGet[ChanM](idx.f.table(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	doc := hmeili.NewMapDocument(m.ID, m.AutoID, m.UpdatedAt.UnixMilli())
	doc["parent"] = m.Parent
	doc["name"] = m.Name
	doc["icon"] = m.Icon
	doc["seq"] = m.Seq
	if len(m.Ctrl) > 0 {
		doc["ctrl"] = m.Ctrl
	}
	if len(m.Tag) > 0 {
		doc["tag"] = m.Tag
	}
	if len(m.Meta) > 0 {
		doc["meta"] = m.Meta
	}
	return doc, nil
}
