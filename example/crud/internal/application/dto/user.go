// Author: Daniel TAN
// Date: 2021-08-18 23:45:12
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:25:21
// FilePath: /trinity-micro/example/crud/internal/application/dto/user.go
// Description:
/*
 * @Author: your name
 * @Date: 2021-08-18 23:45:12
 * @LastEditTime: 2021-09-07 10:46:25
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /trinity-micro/example/internal/application/dto/user.go
 */
package dto

import "github.com/PolarPanda611/trinity-micro/example/crud/internal/application/model"

type GetUserByIDRequest struct {
	*CommonRequest
	CurrentUserID uint64 `header_param:"current_user_id"`
	ID            uint64 `path_param:"id"`
}

type GetUserByIDResponse UserDTO

type ListUserRequest struct {
	*CommonRequest
	Username      string `query_param:"username"`
	CurrentUserID uint64 `header_param:"current_user_id"`
}
type ListUserResponse []UserDTO
type UserDTO struct {
	ID       uint64 `json:"id,string"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func NewListUserResponse(m []model.User) ListUserResponse {
	res := make([]UserDTO, len(m))
	for i, v := range m {
		res[i] = *NewUserDTO(&v)
	}
	return res
}

func NewUserDTO(m *model.User) *UserDTO {
	return &UserDTO{
		ID:       m.ID,
		Username: m.Username,
		Age:      int(m.Age),
	}
}

func NewGetUserByIDResponse(m *model.User) *GetUserByIDResponse {
	d := NewUserDTO(m)
	res := GetUserByIDResponse(*d)
	return &res
}
