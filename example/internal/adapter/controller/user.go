package controller

import (
	"context"
	"net/http"
	"sync"
	"trinity-micro/core/ioc/container"
	"trinity-micro/example/internal/domain/dto"
	"trinity-micro/example/internal/domain/service"
)

func init() {
	UserControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(userControllerImpl)
		},
	}
	container.RegisterInstance("UserController", UserControllerPool)
	container.RegisterController("/test", "UserController",
		container.NewRequestMapping("GET", "/user", "ListUser"),
		container.NewRequestMapping("GET", "/user/{id}", "GetUserByID"),
		container.NewRawRequestMapping("GET", "/testuser", "TestUser"),
	)
}

type userControllerImpl struct {
	UserSrv service.UserService `container:"autowire:true;resource:UserService"`
}

func (c *userControllerImpl) GetUserByID(ctx context.Context, Args struct {
	ID uint64 `path_param:"id"`
}) (*dto.GetUserResponse, error) {
	return c.UserSrv.GetUserID(ctx, uint64(Args.ID))
}

func (c *userControllerImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) ([]dto.GetUserResponse, error) {
	return c.UserSrv.ListUser(ctx, req)
}

func (c *userControllerImpl) TestUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("haha"))
}
