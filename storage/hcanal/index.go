package hcanal

import (
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"strings"
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
	return []string{"INSERT", "UPDATE", "DELETE"}
}

func (h *IndexSyncHandler) OnAlter(alter *Alter) (err error) {
	var curAutoID int64
	defer hlog.Elapse("canal.idx.sync",
		hlog.F(zap.String("table", alter.Table), zap.String("action", alter.Action)),
		func() []zap.Field {
			arr := []zap.Field{zap.Int64("lstAutoID", curAutoID)}
			if err != nil {
				arr = append(arr, zap.Error(err))
			}
			return arr
		})()
	if len(alter.Entities) == 0 {
		return nil
	}
	if strings.EqualFold(alter.Action, "DELETE") {
		var autoIDs []int64
		for _, entity := range alter.Entities {
			autoIDs = append(autoIDs, entity.AutoID)
			curAutoID = entity.AutoID
		}
		err = hmeili.DelDocuments(h.meili, h.indexer, autoIDs)
		if err != nil {
			return err
		}
		return nil
	}

	var docArr []hmeili.Document
	for _, entity := range alter.Entities {
		doc, err := h.indexer.Load(entity.AutoID)
		if err != nil {
			return err
		}
		docArr = append(docArr, doc)
		curAutoID = entity.AutoID
	}
	err = hmeili.AddDocuments(h.meili, h.indexer, docArr)
	if err != nil {
		return err
	}
	return nil
}

func (h *IndexSyncHandler) OnDrop(table string) (err error) {
	defer hlog.Elapse("canal.idx.sync.drop",
		hlog.F(zap.String("table", table)),
		func() []zap.Field {
			if err != nil {
				return []zap.Field{zap.Error(err)}
			}
			return nil
		})()
	err = hmeili.DropIndex(h.meili, h.indexer)
	if err != nil {
		return err
	}
	return nil
}
