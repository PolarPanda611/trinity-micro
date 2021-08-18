package controller

import (
	"encoding/json"
	"net/http"
	"trinity-micro/example/internal/domain/dto"
	"trinity-micro/example/internal/domain/service"
)

var _ UserController = new(userControllerImpl)

type UserController interface {
	GetUserByID(w http.ResponseWriter, r *http.Request)
	ListUser(w http.ResponseWriter, r *http.Request)
}
type userControllerImpl struct {
	userSrv service.UserService
}

func (c *userControllerImpl) GetUserByID(w http.ResponseWriter, r *http.Request) {
	res, err := c.userSrv.GetUserID(r.Context(), 1)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}

func (c *userControllerImpl) ListUser(w http.ResponseWriter, r *http.Request) {
	res, err := c.userSrv.ListUser(r.Context(), &dto.ListUserRequest{
		Username: "123",
	})
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
