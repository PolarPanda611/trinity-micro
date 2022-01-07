package dbx

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64 `gorm:"primarykey;autoIncrement:false" example:"1479429646645936128"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	Version   string     `gorm:"index"`
}

func WithTenant(tenant string, m interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(m)
		return db.Table(tenant + "." + stmt.Schema.Table)
	}
}
