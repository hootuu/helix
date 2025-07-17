package hmeili

import (
	"fmt"
	"github.com/hootuu/helix/components/honce"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cast"
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
		theTask, err := index.AddDocuments(docArr, "id")
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
		gMeiliLogger.Info("added", zap.Int64("task", task.TaskUID))
	}(task)
	return nil
}

func buildInList(ids []int64) string {
	str := ""
	for i, id := range ids {
		if i > 0 {
			str += ", "
		}
		str += `"` + cast.ToString(id) + `"`
	}
	return str
}

func DelDocuments(meili *Meili, indexer Indexer, autoIDs []int64) error {
	if len(autoIDs) == 0 {
		return nil
	}
	filter := []string{"auto_id IN [" + buildInList(autoIDs) + "]"}

	index := meili.Meili().Index(indexer.GetName())
	var task *meilisearch.TaskInfo = nil
	hretry.Universal(func() error {
		theTask, err := index.DeleteDocuments(filter)
		if err != nil {
			hlog.Err("hmeili.DelDocuments", zap.Error(err))
			return err
		}
		task = theTask
		return nil
	})
	if task == nil {
		hlog.Err("hmeili.DelDocuments[final]: Universal err")
		gMeiliLogger.Error("del_err", zap.Any("autoIDs", autoIDs))
		return nil
	}
	go func(task *meilisearch.TaskInfo) {
		status, err := meili.Meili().WaitForTask(task.TaskUID, 1000*time.Millisecond)
		if err != nil || status.Status != meilisearch.TaskStatusSucceeded {
			gMeiliLogger.Error("del_err", zap.Any("autoIDs", autoIDs),
				zap.Int64("task", task.TaskUID),
				zap.Error(err),
				zap.Any("", status),
			)
			return
		}
		gMeiliLogger.Info("deled", zap.Int64("task", task.TaskUID))
	}(task)
	return nil
}

func DropIndex(meili *Meili, indexer Indexer) error {
	var task *meilisearch.TaskInfo = nil
	hretry.Universal(func() error {
		theTask, err := meili.Meili().DeleteIndex(indexer.GetName())
		if err != nil {
			hlog.Err("hmeili.DropIndex", zap.Error(err))
			return err
		}
		task = theTask
		return nil
	})
	if task == nil {
		hlog.Err("hmeili.DropIndex[final]: Universal err")
		gMeiliLogger.Error("drop_err", zap.String("index", indexer.GetName()))
		return nil
	}
	go func(task *meilisearch.TaskInfo) {
		status, err := meili.Meili().WaitForTask(task.TaskUID, 1000*time.Millisecond)
		if err != nil || status.Status != meilisearch.TaskStatusSucceeded {
			gMeiliLogger.Error("drop_err", zap.String("index", indexer.GetName()),
				zap.Int64("task", task.TaskUID),
				zap.Error(err),
				zap.Any("", status),
			)
			return
		}
		gMeiliLogger.Info("dropped", zap.Int64("task", task.TaskUID))
	}(task)
	return nil
}
