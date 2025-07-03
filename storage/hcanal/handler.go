package hcanal

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/hootuu/hyle/hretry"
	"time"
)

func (h *Canal) registerAlterHandler(alter ...AlterHandler) {
	h.alterHandlerArr = append(h.alterHandlerArr, alter...)
}

func (h *Canal) OnRow(event *canal.RowsEvent) error {
	if len(event.Table.Columns) == 0 || len(event.Rows) == 0 {
		return nil
	}
	if len(h.alterHandlerArr) == 0 {
		return nil
	}
	var focusHandler []AlterHandler
	for _, handler := range h.alterHandlerArr {
		if b := h.focusRowsEvent(handler, event); !b {
			return nil
		}
		focusHandler = append(focusHandler, handler)
	}
	if len(focusHandler) == 0 {
		return nil
	}
	alter := &Alter{
		Schema:   event.Table.Schema,
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
			entity.AutoID = row[autoIdIdx].(uint64)
		}
		if len(row) > timestampIdx && row[timestampIdx] != nil {
			entity.Timestamp = row[timestampIdx].(time.Time)
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
	fSchema := handler.Schema()
	fmt.Println("focus", fSchema)
	if fSchema != event.Table.Schema {
		return false
	}
	fmt.Println("focust: ", event.Table.Schema)
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
			if action == event.Action {
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
