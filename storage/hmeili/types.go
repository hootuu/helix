package hmeili

import "github.com/meilisearch/meilisearch-go"

type Doc interface {
	IndexName() string
	GetAutoID() uint64
	GetTimestamp() int64
}

type DocBasic struct {
	AutoID    uint64 `json:"auto_id"`
	Timestamp int64  `json:"timestamp"`
}

func (m *DocBasic) GetAutoID() uint64 {
	return m.AutoID
}

func (m *DocBasic) GetTimestamp() int64 {
	return m.Timestamp
}

type SearchRequest = meilisearch.SearchRequest
