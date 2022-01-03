// Author: Daniel TAN
// Date: 2021-08-18 23:47:20
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 01:26:04
// FilePath: /trinity-micro/example/crud/internal/application/service/user.go
// Description:
package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/core/e"
	"github.com/PolarPanda611/trinity-micro/example/crud/config"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/dto"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/repository"
	"gorm.io/gorm"
)

func init() {
	trinity.RegisterInstance("UserService", &sync.Pool{
		New: func() interface{} {
			return new(userServiceImpl)
		},
	})
}

type UserService interface {
	GetUserID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error)
	ListUser(ctx context.Context, req *dto.ListUserRequest) (*dto.ListUserResponse, error)
}

type userServiceImpl struct {
	UserRepo repository.UserRepository `container:"autowire:true;resource:UserRepository"`
}

func (s *userServiceImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) (*dto.ListUserResponse, error) {
	if req.PageNum != nil {
		var i uint = 1
		req.PageNum = &i
	}
	if req.PageSize != nil {
		req.PageSize = &config.Conf.Application.PageSize
	}
	users, err := s.UserRepo.ListUser(ctx, req.Tenant)
	if err != nil {
		return nil, err
	}
	userCount, err := s.UserRepo.CountUser(ctx, req.Tenant)
	if err != nil {
		return nil, err
	}
	return dto.NewListUserResponse(users, *req.PageSize, *req.PageNum, userCount), nil
}

func (s *userServiceImpl) GetUserID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error) {
	user, err := s.UserRepo.GetUserByID(ctx, req.Tenant, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewError(e.Info, e.ErrRecordNotFound, fmt.Sprintf("recource %v not found", req.ID), err)
		}
		return nil, e.NewError(e.Error, e.ErrExecuteSQL, "excute GetUserByID failed ", err)
	}
	return dto.NewGetUserByIDResponse(user), nil
}
