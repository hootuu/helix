package hcanal

import (
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

type IndexSyncHandler struct {
	table   string
	indexer hmeili.Indexer
	meili   *hmeili.Meili
}

func NewIndexHandler(
	table string,
	indexer hmeili.Indexer,
	meili *hmeili.Meili,
) *IndexSyncHandler {
	return &IndexSyncHandler{
		table:   table,
		indexer: indexer,
		meili:   meili,
	}
}

func (h *IndexSyncHandler) Table() []string {
	return []string{h.table}
}

func (h *IndexSyncHandler) Action() []string {
	return []string{"INSERT", "UPDATE"}
}

func (h *IndexSyncHandler) OnAlter(alter *Alter) (err error) {
	var curAutoID int64
	defer hlog.Elapse("canal.idx.sync",
		hlog.F(zap.String("table", alter.Table), zap.String("action", alter.Action)),
		hlog.E(err, zap.Int64("curAutoID", curAutoID)))()
	if len(alter.Entities) == 0 {
		return nil
	}
	var docArr []hmeili.Document
	for _, entity := range alter.Entities {
		doc, err := h.indexer.Load(entity.AutoID)
		if err != nil {
			return err
		}
		docArr = append(docArr, doc)
	}
	err = hmeili.AddDocuments(h.meili, h.indexer, docArr)
	if err != nil {
		return err
	}
	return nil
}
