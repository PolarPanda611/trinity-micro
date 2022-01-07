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
	"errors"

	"github.com/PolarPanda611/trinity-micro/core/httpx"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/model"
)

type GetUserByIDRequest struct {
	*TenantRequest
	CurrentUserID uint64 `header_param:"current_user_id"`
	ID            uint64 `path_param:"id"`
}

type CreateUserRequest struct {
	*TenantRequest
	CurrentUserID uint64     `header_param:"current_user_id"`
	NewUser       NewUserDTO `body_param:""`
}

func (r *CreateUserRequest) Validate() error {
	if r.NewUser.Email == "" {
		return errors.New("email cannot be empty")
	}
	if r.NewUser.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(r.NewUser.Password) < 8 {
		return errors.New("password cannot be less then 8")
	}
	return nil
}

func (r *CreateUserRequest) Parse() *model.User {
	return &model.User{
		Username: r.NewUser.Username,
		Password: r.NewUser.Password,
		Email:    r.NewUser.Email,
	}
}

type NewUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type UserInfoResponse UserInfoDTO

type ListUserRequest struct {
	*TenantRequest
	*PageRequest
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
	Data []UserInfoDTO
	*httpx.PaginationDTO
}

type UserInfoDTO struct {
	ID       int64  `json:"id,string" example:"1479429646645936128"`
	Username string `json:"username" example:"Daniel"`
	Email    string `json:"email"  example:"daniel@trinity.com"`
	Age      int    `json:"age" example:"18"`
	Gender   string `json:"gender" enums:"male,female" example:"male"`
}

func NewListUserResponse(m []model.User, pageSize, pageNum *uint, total int64) *ListUserResponse {
	res := make([]UserInfoDTO, len(m))
	for i, v := range m {
		res[i] = *NewUserInfoDTO(&v)
	}
	return &ListUserResponse{
		Data:          res,
		PaginationDTO: httpx.NewPaginationDTO(*pageSize, *pageNum, total),
	}
}

func NewUserInfoDTO(m *model.User) *UserInfoDTO {
	return &UserInfoDTO{
		ID:       m.ID,
		Username: m.Username,
		Email:    m.Email,
		Age:      int(m.Age),
	}
}

func NewUserInfoResponse(m *model.User) *UserInfoResponse {
	d := NewUserInfoDTO(m)
	res := UserInfoResponse(*d)
	return &res
}
