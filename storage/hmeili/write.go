package hmeili

import (
	"fmt"
	"github.com/hootuu/helix/components/honce"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
	"time"
)

func InitIndexer(meili *Meili, indexer Indexer) error {
	onceDoCode := fmt.Sprintf("helix_hmeili_index_init_%s_%s", indexer.GetName(), indexer.GetVersion())
	err := honce.Do(onceDoCode, func() error {
		idxMng := meili.Meili().Index(indexer.GetName())
		err := indexer.Setting(idxMng)
		if err != nil {
			hlog.Err("hmeili.InitIndexer",
				zap.String("index", indexer.GetName()),
				zap.String("onceDoCode", onceDoCode),
				zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

var gMeiliLogger = hlog.GetLogger("hmeili")

func AddDocuments(meili *Meili, indexer Indexer, docArr []Document) error {
	index := meili.Meili().Index(indexer.GetName())
	var task *meilisearch.TaskInfo = nil
	hretry.Universal(func() error {
		theTask, err := index.AddDocuments(docArr, "auto_id", "id")
		if err != nil {
			hlog.Err("hmeili.AddDocuments", zap.Error(err))
			return err
		}
		task = theTask
		return nil
	})
	if task == nil {
		hlog.Err("hmeili.AddDocuments[final]: Universal err")
		gMeiliLogger.Error("add_err", zap.Any("docArr", docArr))
		return nil
	}
	go func(task *meilisearch.TaskInfo) {
		status, err := meili.Meili().WaitForTask(task.TaskUID, 1000*time.Millisecond)
		if err != nil || status.Status != meilisearch.TaskStatusSucceeded {
			gMeiliLogger.Error("add_err", zap.Any("docArr", docArr),
				zap.Int64("task", task.TaskUID),
				zap.Error(err),
				zap.Any("", status),
			)
			return
		}
		gMeiliLogger.Info("added", zap.Any("docArr", docArr),
			zap.Int64("task", task.TaskUID))
	}(task)
	return nil
}
