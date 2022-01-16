package dbx

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var (
	_registeredModel = []interface{}{}
	_registeredValue = [][]interface{}{}
)

func RegisterModel(m interface{}, initialValue ...interface{}) {
	_registeredModel = append(_registeredModel, m)
	_registeredValue = append(_registeredValue, initialValue)
}

func PGMigrate(ctx context.Context, tenants ...string) error {
	sessionDB := DB.Session(&gorm.Session{
		NewDB:   true,
		Context: ctx,
	})

	for _, tenant := range tenants {
		for i, v := range _registeredModel {
			if err := sessionDB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS \"%v\" ;", tenant)).Error; err != nil {
				fmt.Println(err)
			}
			if err := sessionDB.Scopes(WithTenant(tenant, v)).AutoMigrate(v); err != nil {
				fmt.Println(err)
			}
			for _, value := range _registeredValue[i] {
				if err := sessionDB.Scopes(WithTenant(tenant, v)).FirstOrCreate(value); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
}

func MysqlMigrate(ctx context.Context, tenants ...string) error {
	sessionDB := DB.Session(&gorm.Session{
		NewDB:   true,
		Context: ctx,
	})

	for _, tenant := range tenants {
		for i, v := range _registeredModel {
			if err := sessionDB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS `%v` ;", tenant)).Error; err != nil {
				fmt.Println(err)
			}
			if err := sessionDB.Scopes(WithTenant(tenant, v)).AutoMigrate(v); err != nil {
				fmt.Println(err)
			}
			for _, value := range _registeredValue[i] {
				if err := sessionDB.Scopes(WithTenant(tenant, v)).FirstOrCreate(value); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
}
