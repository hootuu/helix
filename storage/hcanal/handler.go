package hcanal

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyle/hsys"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

func (h *Canal) OnDDL(_ *replication.EventHeader, _ mysql.Position, queryEvent *replication.QueryEvent) error {
	if hsys.RunMode().IsProd() || hsys.RunMode().IsPre() {
		return nil
	}
	if len(h.alterHandlerArr) == 0 {
		return nil
	}
	queryStr := strings.ToLower(string(queryEvent.Query))
	b, arr := analyseDropTableDDL(queryStr)
	if !b {
		return nil
	}
	if len(arr) == 0 {
		return nil
	}
	for _, table := range arr {
		for _, handler := range h.alterHandlerArr {
			focusTblArr := handler.Table()
			if len(focusTblArr) == 0 {
				continue
			}
			for _, tbl := range focusTblArr {
				if strings.EqualFold(tbl, table) {
					err := handler.OnDrop(table)
					if err != nil {
						hlog.Fix("helix.canal.OnDDL", zap.Error(err))
						continue //USE CONTINUE
					}
				}
			}
		}
	}
	return nil
}

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

func analyseDropTableDDL(query string) (bool, []string) {
	hQuery := strings.ToLower(strings.TrimSpace(query))
	reMultiLine := regexp.MustCompile(`/\*[^*]*\*+(?:[^/*][^*]*\*+)*/`)
	hQuery = reMultiLine.ReplaceAllString(hQuery, "")
	reSingleLine := regexp.MustCompile(`(?:--|#).*`)
	hQuery = reSingleLine.ReplaceAllString(hQuery, "")
	if !strings.HasPrefix(hQuery, "drop table") {
		return false, nil
	}

	tableSpec := strings.TrimPrefix(hQuery, "drop table")
	tableSpec = strings.TrimPrefix(tableSpec, "if exists")
	tableSpec = strings.TrimSpace(tableSpec)

	tables := strings.Split(tableSpec, ",")
	var result []string

	for _, t := range tables {
		clean := strings.TrimSpace(t)
		clean = strings.Trim(clean, ";")
		clean = strings.Trim(clean, "`")

		if clean != "" {
			result = append(result, clean)
		}
	}

	return true, result
}
