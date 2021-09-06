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
	GetUserID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error)
	ListUser(ctx context.Context, req *dto.ListUserRequest) (dto.ListUserResponse, error)
}

type userServiceImpl struct {
	UserRepo repository.UserRepository `container:"autowire:true;resource:UserRepository"`
}

func (s *userServiceImpl) GetUserID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error) {
	user, err := s.UserRepo.GetUserByID(ctx, req.CurrentUserID, req.ID)
	if err != nil {
		return nil, err
	}
	return dto.NewGetUserByIDResponse(user), nil
}

func (s *userServiceImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) (dto.ListUserResponse, error) {
	users, err := s.UserRepo.ListUser(ctx, req.CurrentUserID)
	if err != nil {
		return nil, err
	}
	return dto.NewListUserResponse(users), nil
}
