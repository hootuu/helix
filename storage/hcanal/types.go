package hcanal

import "time"

type Entity struct {
	AutoID    uint64    `json:"auto_id"`
	Timestamp time.Time `json:"timestamp"`
}
type Alter struct {
	Schema   string    `json:"schema"`
	Table    string    `json:"table"`
	Action   string    `json:"action"`
	Entities []*Entity `json:"entities"`
}

type AlterHandler interface {
	Schema() string
	Table() []string
	Action() []string
	OnAlter(alter *Alter) error
}
