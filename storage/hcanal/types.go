package hcanal

type Entity struct {
	AutoID    int64 `json:"auto_id"`
	Timestamp int64 `json:"timestamp"`
}
type Alter struct {
	Table    string    `json:"table"`
	Action   string    `json:"action"`
	Entities []*Entity `json:"entities"`
}

type AlterHandler interface {
	Table() []string
	Action() []string
	OnAlter(alter *Alter) error
	OnDrop(table string) error
}
