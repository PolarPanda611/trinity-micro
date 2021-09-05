package controller

import (
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
	)
}

type userControllerImpl struct {
	// userSrv service.UserService `container:"autowire:true"`
}

func (c *userControllerImpl) GetUserByID(Args struct {
	ID int `path_param:"id"`
}) int {
	return Args.ID
}

func (c *userControllerImpl) ListUser() string {
	return "haha"
}
