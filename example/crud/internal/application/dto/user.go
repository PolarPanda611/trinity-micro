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

import (
	"github.com/PolarPanda611/trinity-micro/core/httpx"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/model"
)

type GetUserByIDRequest struct {
	*CommonRequest
	CurrentUserID uint64 `header_param:"current_user_id"`
	ID            uint64 `path_param:"id"`
}

type GetUserByIDResponse UserDTO

type ListUserRequest struct {
	*CommonRequest
	PageSize      *uint   `query_param:"pageSize" validate:"required"`
	PageNum       *uint   `query_param:"current"`
	UsernameIlike *string `query_param:"username__ilike"`
	Age           *int    `query_param:"age"`
	CurrentUserID uint64  `header_param:"current_user_id"`
}

type ListUserPageQuery struct {
	PageSize *uint
	PageNum  *uint
	*ListUserQuery
}
type ListUserQuery struct {
	UsernameIlike *string
	Age           *int
	CurrentUserID uint64
}

func (r *ListUserRequest) ParseQuery() *ListUserQuery {
	return &ListUserQuery{
		UsernameIlike: r.UsernameIlike,
		Age:           r.Age,
		CurrentUserID: r.CurrentUserID,
	}
}
func (r *ListUserRequest) ParsePageQuery() *ListUserPageQuery {
	return &ListUserPageQuery{
		PageSize:      r.PageSize,
		PageNum:       r.PageNum,
		ListUserQuery: r.ParseQuery(),
	}
}

type ListUserResponse struct {
	Data []UserDTO
	*httpx.PaginationDTO
}

type UserDTO struct {
	ID       uint64 `json:"id,string" example:"1479429646645936128"`
	Username string `json:"username" example:"Daniel"`
	Age      int    `json:"age" example:"18"`
	Gender   string `json:"gender" enums:"male,female" example:"male"`
}

func NewListUserResponse(m []model.User, pageSize, pageNum *uint, total int64) *ListUserResponse {
	res := make([]UserDTO, len(m))
	for i, v := range m {
		res[i] = *NewUserDTO(&v)
	}
	return &ListUserResponse{
		Data:          res,
		PaginationDTO: httpx.NewPaginationDTO(*pageSize, *pageNum, total),
	}
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
