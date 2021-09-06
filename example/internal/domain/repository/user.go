package repository

import (
	"context"
	"fmt"
	"sync"
	"trinity-micro/core/ioc/container"
	"trinity-micro/example/internal/domain/model"
)

func init() {
	container.RegisterInstance("UserRepository", &sync.Pool{
		New: func() interface{} {
			return new(userRepositoryImpl)
		},
	})
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id uint64) (*model.User, error)
	ListUser(ctx context.Context) ([]model.User, error)
}

type userRepositoryImpl struct {
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	return &model.User{
		ID:       id,
		Username: fmt.Sprintf("daniel_usernname_%v", id),
		Password: fmt.Sprintf("daniel_password_%v", id),
		Email:    fmt.Sprintf("daniel_email_%v", id),
		Age:      1,
		Gender:   1,
	}, nil
}

func (r *userRepositoryImpl) ListUser(ctx context.Context) ([]model.User, error) {
	var len uint = 6
	list := make([]model.User, len)
	for len > 0 {
		list[len-1] = model.User{
			ID:       uint64(len),
			Username: fmt.Sprintf("daniel_usernname_%v", len),
			Password: fmt.Sprintf("daniel_password_%v", len),
			Email:    fmt.Sprintf("daniel_email_%v", len),
			Age:      len,
			Gender:   len,
		}
		len--
	}
	return list, nil
}
