package httpx

import "math"

type PaginationDTO struct {
	Total     int64 `json:"total" example:"120"`
	Current   uint  `json:"current" example:"3"`
	TotalPage uint  `json:"total_page" example:"6"`
	PageSize  uint  `json:"page_size" example:"20"`
}

func NewPaginationDTO(pageSize, pageNum uint, total int64) *PaginationDTO {
	return &PaginationDTO{
		Total:     total,
		PageSize:  pageSize,
		TotalPage: uint(math.Ceil(float64(total) / float64(pageSize))),
		Current:   pageNum,
	}
}
