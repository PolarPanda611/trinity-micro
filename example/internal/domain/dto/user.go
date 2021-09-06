package dto

import "trinity-micro/example/internal/domain/model"

type GetUserByIDRequest struct {
	CurrentUserID uint64 `header_param:"current_user_id"`
	ID            uint64 `path_param:"id"`
}

type GetUserByIDResponse UserDTO

type ListUserRequest struct {
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
