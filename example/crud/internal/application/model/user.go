// Author: Daniel TAN
// Date: 2021-10-02 01:20:48
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:57:40
// FilePath: /trinity-micro/example/crud/internal/application/model/user.go
// Description:
package model

type User struct {
	ID       uint64
	Username string
	Password string
	Email    string
	Age      uint
	Gender   uint
	// Orders    []Order
	CreatedBy uint64
}
