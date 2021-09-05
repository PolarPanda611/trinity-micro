package controller

import (
	"net/http"
	"sync"
	"trinity-micro/core/ioc/container"
)

var _ UserController = new(userControllerImpl)

func init() {
	UserControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(userControllerImpl)
		},
	}
	container.RegisterInstance("UserController", UserControllerPool)
	container.RegisterController("/test", "UserController",
		container.NewRequestMapping("GET", "/user", "ListUser"),
	)
}

type UserController interface {
	GetUserByID(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
}
type userControllerImpl struct {
	// userSrv service.UserService `container:"autowire:true"`
}

func (c *userControllerImpl) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// res, err := c.userSrv.GetUserID(r.Context(), 1)
	// if err != nil {
	// 	w.WriteHeader(400)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	w.WriteHeader(200)
	w.Write([]byte("pong"))
}

func (c *userControllerImpl) ListUser(w http.ResponseWriter, r *http.Request) {
	// res, err := c.userSrv.ListUser(r.Context(), &dto.ListUserRequest{
	// 	Username: "123",
	// })
	// if err != nil {
	// 	w.WriteHeader(400)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// data, err := json.Marshal(res)
	// if err != nil {
	// 	w.WriteHeader(400)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// w.WriteHeader(200)
	// w.Write(data)
}
