package dto

import "trinity-micro/example/internal/domain/model"

type ListUserRequest struct {
	Username      string `query_param:"username"`
	Age           int    `query_param:"age"`
	CurrentUserID uint64 `header_param:"current_user_id"`
}

type GetUserResponse struct {
	ID       uint64 `json:"id,string"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func NewGetUserResponse(m *model.User) *GetUserResponse {
	return &GetUserResponse{
		ID:       m.ID,
		Username: m.Username,
		Age:      int(m.Age),
	}
}
