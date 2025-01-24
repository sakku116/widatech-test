package dto

type BaseJSONResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
	Data    interface{} `json:"data"`
}

type BasePaginationRespData struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
}

func (dto *BasePaginationRespData) SetPagination(total int64, page int, limit int) {
	dto.Total = total
	dto.CurrentPage = page
	if limit > 0 && page > 0 {
		dto.TotalPage = (total + int64(limit) - 1) / int64(limit)
	}
}
