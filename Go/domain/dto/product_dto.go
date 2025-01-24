package dto

import (
	"backend/domain/enum"
	"backend/domain/model"
	"fmt"
)

type ProductRepo_GetListParams struct {
	Query     *string
	QueryBy   *string // leave empty to query by all
	Page      *int
	Limit     *int
	SortOrder *enum.SortOrder
	SortBy    *string
	DoCount   bool
}

func (params *ProductRepo_GetListParams) Validate() error {

	if params.SortOrder != nil && !(*params.SortOrder).IsValid() {
		return fmt.Errorf("invalid sort_order")
	}

	tmp := model.Product{}
	if params.QueryBy != nil {
		queriableFields := tmp.GetProps().QueriableFields
		contain := false
		for _, field := range queriableFields {
			if *params.QueryBy == field {
				contain = true
				break
			}
		}
		if !contain {
			return fmt.Errorf("invalid query_by")
		}
	}

	if params.SortBy != nil {
		sortableFields := tmp.GetProps().SortableFields
		contain := false
		for _, field := range sortableFields {
			if *params.SortBy == field {
				contain = true
				break
			}
		}
		if !contain {
			return fmt.Errorf("invalid sort_by")
		}
	}

	return nil
}
