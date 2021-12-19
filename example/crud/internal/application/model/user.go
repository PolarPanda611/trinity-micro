// Author: Daniel TAN
// Date: 2021-10-02 01:20:48
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:57:40
// FilePath: /trinity-micro/example/crud/internal/application/model/user.go
// Description:
package model

import "github.com/PolarPanda611/trinity-micro/core/dbx"

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
