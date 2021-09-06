package controller

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"trinity-micro/core/ioc/container"
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
	// userSrv service.UserService `container:"autowire:true"`
}

func (c *userControllerImpl) GetUserByID(ctx context.Context, Args struct {
	ID int `path_param:"id"`
}) int {
	return Args.ID
}

func (c *userControllerImpl) ListUser(ctx context.Context, Args struct {
	Query url.Values `query_param:""`
}) url.Values {
	return Args.Query
}

func (c *userControllerImpl) TestUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("haha"))
}
