package dbx

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var (
	_registeredModel = []interface{}{}
)

func RegisterModel(m interface{}) {
	_registeredModel = append(_registeredModel, m)
}

func Migrate(ctx context.Context, tenants ...string) error {
	sessionDB := DB.Session(&gorm.Session{
		NewDB:   true,
		Context: ctx,
	})

	for _, tenant := range tenants {
		for _, v := range _registeredModel {
			if err := sessionDB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS \"%v\" ;", tenant)).Error; err != nil {
				return err
			}
			if err := sessionDB.Scopes(WithTenant(tenant, v)).AutoMigrate(v); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
