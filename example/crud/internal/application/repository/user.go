// Author: Daniel TAN
// Date: 2021-08-19 00:01:37
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:53:55
// FilePath: /trinity-micro/example/crud/internal/application/repository/user.go
// Description:
package repository

import (
	"context"
	"sync"

	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/model"

	"github.com/PolarPanda611/trinity-micro/core/dbx"
)

func init() {
	trinity.RegisterInstance("UserRepository", &sync.Pool{
		New: func() interface{} {
			return new(userRepositoryImpl)
		},
	})
}

var _ UserRepository = new(userRepositoryImpl)

type UserRepository interface {
	GetUserByID(ctx context.Context, tenant string, ID uint64) (*model.User, error)
	ListUser(ctx context.Context, tenant string) ([]model.User, error)
	CountUser(ctx context.Context, tenant string) (int64, error)
}

type userRepositoryImpl struct {
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, tenant string, ID uint64) (*model.User, error) {
	res := &model.User{}
	if err := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	).
		Where("id = ?", ID).First(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *userRepositoryImpl) ListUser(ctx context.Context, tenant string) ([]model.User, error) {
	res := []model.User{}
	if err := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *userRepositoryImpl) CountUser(ctx context.Context, tenant string) (int64, error) {
	var c int64
	if err := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	).
		Count(&c).Error; err != nil {
		return 0, err
	}
	return c, nil
}
