package dbx

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func WithTenant(tenant string, m interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(m)
		return db.Table(tenant + "." + stmt.Schema.Table)
	}
}
