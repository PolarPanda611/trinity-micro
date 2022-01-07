// Author: Daniel TAN
// Date: 2021-08-19 00:01:37
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:53:55
// FilePath: /trinity-micro/example/crud/internal/application/repository/user.go
// Description:
package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/dto"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/model"
	"gorm.io/gorm"

	"github.com/PolarPanda611/trinity-micro/core/dbx"
	"github.com/PolarPanda611/trinity-micro/core/e"
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
	ListUser(ctx context.Context, tenant string, query *dto.ListUserPageQuery) ([]model.User, error)
	CountUser(ctx context.Context, tenant string, query *dto.ListUserQuery) (int64, error)
	CreateUser(ctx context.Context, tenant string, newUser *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, tenant string, id uint64, version string, change map[string]interface{}) error
}

type userRepositoryImpl struct {
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, tenant string, ID uint64) (*model.User, error) {
	res := &model.User{}
	if err := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	).
		Where("id = ?", ID).First(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewError(e.Info, e.ErrRecordNotFound, fmt.Sprintf("userRepositoryImpl.GetUserByID failed tenant: %v id: %v", tenant, ID))
		}
		return nil, e.NewError(e.Error, e.ErrExecuteSQL, fmt.Sprintf("userRepositoryImpl.GetUserByID failed tenant: %v id: %v,err: %v ", tenant, ID, err))
	}
	return res, nil
}

func (r *userRepositoryImpl) ListUser(ctx context.Context, tenant string, query *dto.ListUserPageQuery) ([]model.User, error) {
	db := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	)
	if query.UsernameIlike != nil {
		db = db.Where("username ilike ?", "%"+*query.UsernameIlike+"%")
	}
	if query.Age != nil {
		db = db.Where("age = ?", query.Age)
	}
	res := []model.User{}
	if err := db.Find(&res).Error; err != nil {
		return nil, e.NewError(e.Error, e.ErrExecuteSQL, fmt.Sprintf("userRepositoryImpl.ListUser failed, tenant: %v, error: %v ", tenant, err))
	}
	return res, nil
}

func (r *userRepositoryImpl) CountUser(ctx context.Context, tenant string, query *dto.ListUserQuery) (int64, error) {
	db := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	)
	if query.UsernameIlike != nil {
		db = db.Where("username ilike ?", "%"+*query.UsernameIlike+"%")
	}
	if query.Age != nil {
		db = db.Where("age = ?", query.Age)
	}
	var c int64
	if err := db.Count(&c).Error; err != nil {
		return 0, e.NewError(e.Error, e.ErrExecuteSQL, fmt.Sprintf("userRepositoryImpl.CountUser failed, tenant: %v, error: %v ", tenant, err))
	}
	return c, nil
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, tenant string, newUser *model.User) (*model.User, error) {
	db := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	)
	if err := db.Create(newUser).Error; err != nil {
		return nil, e.NewError(e.Error, e.ErrExecuteSQL, fmt.Sprintf("userRepositoryImpl.CreateUser failed, tenant: %v, error: %v ", tenant, err))
	}
	return newUser, nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, tenant string, id uint64, version string, change map[string]interface{}) error {
	db := dbx.FromCtx(ctx).Scopes(
		dbx.WithTenant(tenant, &model.User{}),
	)
	res := db.Where("id = ?", id).Where("version = ?", version).Updates(change)
	if err := res.Error; err != nil {
		return e.NewError(e.Error, e.ErrExecuteSQL, fmt.Sprintf("userRepositoryImpl.UpdateUser failed, tenant: %v, id: %v version: %v, error: %v ", tenant, id, version, err))
	}
	if res.RowsAffected != 1 {
		return e.NewError(e.Info, e.ErrDBUpdateZeroLine, fmt.Sprintf("userRepositoryImpl.UpdateUser tenant: %v, id: %v version: %v 0 rows affected", tenant, id, version))
	}
	return nil
}
