package hdb

import (
	"errors"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Find[T any](find func() *gorm.DB) ([]*T, error) {
	var arr []*T
	tx := find().Find(&arr)
	if tx.Error != nil {
		hlog.Err("hdb.Find", zap.Error(tx.Error))
		return nil, errors.New("db error:" + tx.Error.Error())
	}
	return arr, nil
}

func PageFind[T any](page *pagination.Page, find func() *gorm.DB) (*pagination.Pagination[T], error) {
	if page == nil {
		page = pagination.PageNormal()
	}
	var arr []*T
	var count int64
	countTx := find().Count(&count)
	if countTx.Error != nil {
		hlog.Err("hdb.PageFind:Count()", zap.Error(countTx.Error))
		return nil, errors.New("db.count error:" + countTx.Error.Error())
	}
	pageTx := find().Limit(int(page.Size)).Offset(int((page.Numb - 1) * page.Size)).Find(&arr)
	if pageTx.Error != nil {
		hlog.Err("hdb.PageFind:Page()", zap.Error(pageTx.Error))
		return nil, errors.New("db.find error:" + pageTx.Error.Error())
	}
	return pagination.NewPagination[T](pagination.PagingOfPage(page).WithCount(count), arr), nil
}
