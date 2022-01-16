package dbx

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        int64 `gorm:"primarykey;autoIncrement:false" example:"1479429646645936128"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	Version   string     `gorm:"type:varchar(50);index"`
}

func WithTenant(tenant string, m interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(m)
		return db.Table(tenant + "." + stmt.Schema.Table)
	}
}

// WithPagenation
// build the gorm condition with page_num and page_size
// the page_num start from 1
func WithPagenation(pageNum int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit := pageSize
		offset := (pageNum - 1) * pageSize
		return db.Limit(limit).Offset(offset)
	}
}
