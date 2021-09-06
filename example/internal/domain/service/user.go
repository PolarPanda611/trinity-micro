package service

import (
	"context"
	"sync"
	"trinity-micro/core/ioc/container"
	"trinity-micro/example/internal/domain/dto"
	"trinity-micro/example/internal/domain/repository"
)

func init() {
	container.RegisterInstance("UserService", &sync.Pool{
		New: func() interface{} {
			return new(userServiceImpl)
		},
	})
}

type UserService interface {
	GetUserID(ctx context.Context, id uint64) (*dto.GetUserResponse, error)
	ListUser(ctx context.Context, req *dto.ListUserRequest) ([]dto.GetUserResponse, error)
}

type userServiceImpl struct {
	UserRepo repository.UserRepository `container:"autowire:true;resource:UserRepository"`
}

func (s *userServiceImpl) GetUserID(ctx context.Context, id uint64) (*dto.GetUserResponse, error) {
	user, err := s.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.NewGetUserResponse(user), nil
}

func (s *userServiceImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) ([]dto.GetUserResponse, error) {
	users, err := s.UserRepo.ListUser(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.GetUserResponse, len(users))
	for i, v := range users {
		res[i] = *dto.NewGetUserResponse(&v)
	}
	return res, nil
}
