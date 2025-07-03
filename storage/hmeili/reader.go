package hmeili

import (
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func Find(meili *Meili, indexName string, req *SearchRequest, page *pagination.Page) (*pagination.Pagination[any], error) {
	paging := pagination.PagingOfPage(page)
	if req == nil {
		req = &SearchRequest{
			Query: "*",
			Sort:  []string{"auto_id:asc"},
		}
	}
	req.Offset = paging.Skip()
	req.Limit = paging.Limit()

	index := meili.Meili().Index(indexName)
	result, err := index.Search("*", req)
	if err != nil {
		hlog.Err("hmeili.Find", zap.Error(err))
		return nil, err
	}
	//
	//data, err := result.MarshalJSON()
	//if err != nil {
	//	hlog.Err("hmeili.Find", zap.Error(err))
	//	return nil, err
	//}
	//fmt.Println(string(data))
	//
	//var arr []T
	//err = json.Unmarshal(data, &arr)
	//if err != nil {
	//	return nil, err
	//}
	return pagination.NewPagination(paging, result.Hits), nil
}
