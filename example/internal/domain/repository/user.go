package repository

import (
	"context"
	"fmt"
	"sync"
	"trinity-micro/core/e"
	"trinity-micro/core/ioc/container"
	"trinity-micro/example/internal/domain/model"
)

var (
	userLake = []model.User{}
)

func init() {
	container.RegisterInstance("UserRepository", &sync.Pool{
		New: func() interface{} {
			return new(userRepositoryImpl)
		},
	})
	var len uint = 6
	for len > 0 {
		userLake = append(userLake, model.User{
			ID:        uint64(len),
			Username:  fmt.Sprintf("daniel_usernname_%v", len),
			Password:  fmt.Sprintf("daniel_password_%v", len),
			Email:     fmt.Sprintf("daniel_email_%v", len),
			Age:       len,
			Gender:    len,
			CreatedBy: uint64(len/4 + 1),
		})
		len--
	}
}

type UserRepository interface {
	GetUserByID(ctx context.Context, currentUserID uint64, ID uint64) (*model.User, error)
	ListUser(ctx context.Context, currentUserID uint64) ([]model.User, error)
}

type userRepositoryImpl struct {
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, currentUserID uint64, ID uint64) (*model.User, error) {
	for _, v := range userLake {
		// auth check
		if v.CreatedBy != currentUserID {
			continue
		}
		if v.ID == ID {
			return &v, nil
		}

	}
	return nil, e.NewError(e.Info, e.ErrRecordNotFound, fmt.Sprintf("user not found => id: %v", ID))
}

func (r *userRepositoryImpl) ListUser(ctx context.Context, currentUserID uint64) ([]model.User, error) {
	res := []model.User{}
	for _, v := range userLake {
		if v.CreatedBy == currentUserID {
			res = append(res, v)
		}
	}
	return res, nil
}
