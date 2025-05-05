package hnid

import (
	"fmt"
	"github.com/hootuu/hyle/data/hcast"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"strings"
)

type NID struct {
	Timestamp    uint64 `json:"timestamp"`
	Biz          uint64 `json:"biz"`
	AutoInc      uint64 `json:"auto_inc"`
	TimestampLen uint8  `json:"timestamp_len"`
	BizLen       uint8  `json:"biz_len"`
	AutoIncLen   uint8  `json:"auto_inc_len"`
}

func (id NID) reset() {
	id.Timestamp = 0
	id.TimestampLen = 0
	id.AutoInc = 0
	id.AutoIncLen = 0
	id.Biz = 0
	id.BizLen = 0
}

func (id NID) itemToString(len uint8, val uint64) string {
	format := fmt.Sprintf("%%0%dd", len)
	str := fmt.Sprintf(format, val)
	return str
}

func (id NID) ToString() string {
	var idBuf strings.Builder
	if id.BizLen > 0 {
		idBuf.WriteString(id.itemToString(id.BizLen, id.Biz))
	}
	if id.TimestampLen > 0 {
		idBuf.WriteString(fmt.Sprintf("%d%d", id.TimestampLen, id.Timestamp))
	}
	if id.AutoIncLen > 0 {
		idBuf.WriteString(id.itemToString(id.AutoIncLen, id.AutoInc))
	}
	return idBuf.String()
}

func (id NID) ToUint64() uint64 {
	strID := id.ToString()
	if len(strID) > 19 {
		hlog.Err("helix.hnid.nid.ToUint64: len(strID) > 19", zap.String("id", strID))
		return 0
	}
	return hcast.ToUint64(strID)
}
