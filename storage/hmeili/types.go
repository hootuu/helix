package hmeili

import (
	"github.com/meilisearch/meilisearch-go"
)

type SearchRequest = meilisearch.SearchRequest

type Indexer interface {
	GetName() string
	GetVersion() string
	Setting(index meilisearch.IndexManager) error
	Load(autoID int64) (Document, error)
}

type Document interface {
	GetID() any
	GetAutoID() int64
	GetTimestamp() int64
}

type MapDocument map[string]interface{}

func NewMapDocument(id any, autoID int64, timestamp int64) MapDocument {
	doc := make(MapDocument)
	doc["id"] = id
	doc["auto_id"] = autoID
	doc["timestamp"] = timestamp
	return doc
}

func (d MapDocument) GetID() any {
	return d["id"]
}

func (d MapDocument) GetAutoID() int64 {
	obj, ok := d["auto_id"]
	if !ok {
		return 0
	}
	if autoID, ok := obj.(int64); ok {
		return autoID
	}
	return 0
}

func (d MapDocument) GetTimestamp() int64 {
	obj, ok := d["timestamp"]
	if !ok {
		return 0
	}
	if timestamp, ok := obj.(int64); ok {
		return timestamp
	}
	return 0
}

func (d MapDocument) Mix(prefix string, data map[string]interface{}) {
	if len(data) == 0 {
		return
	}
	for k, v := range data {
		mixKey := prefix + "_" + k
		d[mixKey] = v
	}
}
