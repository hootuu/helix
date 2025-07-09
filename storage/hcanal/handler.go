package hcanal

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/spf13/cast"
	"strings"
)

func (h *Canal) OnRow(event *canal.RowsEvent) error {
	if len(event.Table.Columns) == 0 || len(event.Rows) == 0 {
		hlog.Err("hcanal.OnRow: no columns or rows")
		return nil
	}
	if len(h.alterHandlerArr) == 0 {
		hlog.Logger().Warn("hcanal.OnRow: len(h.alterHandlerArr) == 0")
		return nil
	}
	var focusHandler []AlterHandler
	for _, handler := range h.alterHandlerArr {
		if b := h.focusRowsEvent(handler, event); !b {
			continue
		}
		focusHandler = append(focusHandler, handler)
	}
	if len(focusHandler) == 0 {
		//hlog.Info("hcanal.OnRow: len(h.alterHandlerArr) == 0")
		return nil
	}
	alter := &Alter{
		Table:    event.Table.Name,
		Action:   event.Action,
		Entities: nil,
	}
	var autoIdIdx, timestampIdx int
	for i, col := range event.Table.Columns {
		switch col.Name {
		case "auto_id":
			autoIdIdx = i
		case "updated_at":
			timestampIdx = i
		default:
			continue
		}
	}
	for _, row := range event.Rows {
		entity := &Entity{}
		if len(row) > autoIdIdx && row[autoIdIdx] != nil {
			entity.AutoID = row[autoIdIdx].(int64)
		}
		if len(row) > timestampIdx && row[timestampIdx] != nil {
			entity.Timestamp = cast.ToTime(row[timestampIdx]).UnixMilli()
		}
		alter.Entities = append(alter.Entities, entity)
	}
	for _, handler := range focusHandler {
		hretry.Universal(func() error {
			err := handler.OnAlter(alter)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return nil
}

func (h *Canal) String() string {
	return "HelixCanalHandler"
}

func (h *Canal) focusRowsEvent(handler AlterHandler, event *canal.RowsEvent) bool {
	fTable := handler.Table()
	if len(fTable) > 0 {
		beIn := false
		for _, tbl := range fTable {
			if tbl == event.Table.Name {
				beIn = true
				break
			}
		}
		if !beIn {
			return false
		}
	}
	fAction := handler.Action()
	if len(fAction) > 0 {
		beIn := false
		for _, action := range fAction {
			if strings.EqualFold(action, event.Action) {
				beIn = true
				break
			}
		}
		if !beIn {
			return false
		}
	}
	return true
}
