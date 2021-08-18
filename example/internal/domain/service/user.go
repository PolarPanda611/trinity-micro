package service

import (
	"context"
	"trinity-micro/example/internal/domain/dto"
	"trinity-micro/example/internal/domain/repository"
)

type UserService interface {
	GetUserID(ctx context.Context, id uint64) (*dto.GetUserResponse, error)
	ListUser(ctx context.Context, req *dto.ListUserRequest) ([]dto.GetUserResponse, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func (s *userServiceImpl) GetUserID(ctx context.Context, id uint64) (*dto.GetUserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.NewGetUserResponse(user), nil
}

func (s *userServiceImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) ([]dto.GetUserResponse, error) {
	users, err := s.userRepo.ListUser(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.GetUserResponse, len(users))
	for i, v := range users {
		res[i] = *dto.NewGetUserResponse(&v)
	}
	return res, nil
}
