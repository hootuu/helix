package hmeili

import (
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func Search(
	meili *Meili,
	indexName string,
	filter string,
	searchAttr []SearchAttr,
	sort []string,
	page *pagination.Page,
) (*pagination.Pagination[any], error) {
	paging := pagination.PagingOfPage(page)
	req := &SearchRequest{
		Offset:              paging.Skip(),
		Limit:               paging.Limit(),
		Filter:              filter,
		Sort:                sort,
		ShowMatchesPosition: true,
	}
	if len(req.Sort) == 0 {
		req.Sort = []string{"auto_id:desc"}
	}

	index := meili.Meili().Index(indexName)
	query := ""
	var attributesToSearchOn []string
	if len(searchAttr) > 0 {
		for _, attr := range searchAttr {
			if query != "" {
				query += " "
			}
			query += attr.Val
			attributesToSearchOn = append(attributesToSearchOn, attr.Attr)
		}
	}

	if len(attributesToSearchOn) > 0 {
		req.AttributesToSearchOn = attributesToSearchOn
	}

	result, err := index.Search(query, req)

	if err != nil {
		hlog.Err("hmeili.Find", zap.Error(err))
		return nil, err
	}
	paging.WithCount(result.EstimatedTotalHits)
	return pagination.NewPagination(paging, result.Hits), nil
}

func Filter(
	meili *Meili,
	indexName string,
	filter string,
	sort []string,
	page *pagination.Page,
) (*pagination.Pagination[any], error) {
	return Search(meili, indexName, filter, nil, sort, page)
}

func Find(meili *Meili, indexName string, req *SearchRequest, page *pagination.Page) (*pagination.Pagination[any], error) {
	paging := pagination.PagingOfPage(page)
	if req == nil {
		req = &SearchRequest{
			Query: "",
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

	paging.WithCount(result.TotalHits)
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
