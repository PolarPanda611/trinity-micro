// Author: Daniel TAN
// Date: 2021-10-02 01:20:48
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:57:40
// FilePath: /trinity-micro/example/crud/internal/application/model/user.go
// Description:
package model

import (
	"github.com/PolarPanda611/trinity-micro/core/dbx"
	"github.com/PolarPanda611/trinity-micro/core/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func init() {
	dbx.RegisterModel(&User{})
}

type User struct {
	dbx.Model
	Username string
	Password string
	Email    string
	Age      uint
	Gender   uint
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == 0 {
		u.ID = utils.GetSnowflakeID()
	}
	if u.Version == "" {
		u.Version = uuid.NewV4().String()
	}
	return
}
