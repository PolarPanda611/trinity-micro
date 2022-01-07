package dbx

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64 `gorm:"primarykey;autoIncrement:false" example:"1479429646645936128"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

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
func WithTenant(tenant string, m interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(m)
		return db.Table(tenant + "." + stmt.Schema.Table)
	}
}
