// Author: Daniel TAN
// Date: 2021-08-18 23:39:51
// LastEditors: Daniel TAN
// LastEditTime: 2021-09-27 23:11:06
// FilePath: /trinity-micro/example/internal/adapter/controller/user.go
// Description:
package controller

import (
	"context"
	"net/http"
	"sync"

	"github.com/PolarPanda611/trinity-micro/example/internal/application/service"

	"github.com/PolarPanda611/trinity-micro/core/ioc/container"

	"github.com/PolarPanda611/trinity-micro/example/internal/application/dto"

	"github.com/PolarPanda611/trinity-micro/core/httpx"
)

func init() {
	UserControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(userControllerImpl)
		},
	}
	container.RegisterInstance("UserController", UserControllerPool)
	container.RegisterController("/diparam", "UserController",
		container.NewRequestMapping("GET", "/users", "ListUser"),
		container.NewRequestMapping("GET", "/users/{id}", "GetUserByID"),
	)
	container.RegisterController("/raw", "UserController",
		container.NewRequestMapping("GET", "/users", "ListUser"),
		container.NewRequestMapping("GET", "/users/{id}", "GetUserByID"),
	)
}

type userControllerImpl struct {
	UserSrv service.UserService `container:"autowire:true;resource:UserService"`
}

// @Summary: swagger summary
// @Description: swagger description
// @Accept: json
// @Produce: json
// @Param:  id path int true "Account ID"
// @Success: 200 {object} dto.GetUserByIDResponse
// @Failure: 400,404 {object} httputil.HTTPError
// @Router: /accounts/{id} [get]
func (c *userControllerImpl) GetUserByID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error) {
	return c.UserSrv.GetUserID(ctx, req)
}

// @Summary: swagger summary
// @Description: swagger description
// @Accept: json
// @Produce: json
// @Param:  id path int true "Account ID"
// @Success: 200 {object} model.Account
// @Header: 200 {string} Token "qwerty"
// @Failure: 400,404 {object} httputil.HTTPError
// @Router: /accounts/{id} [get]
func (c *userControllerImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) (dto.ListUserResponse, error) {
	return c.UserSrv.ListUser(ctx, req)
}

func (c *userControllerImpl) RawGetUserByID(w http.ResponseWriter, r *http.Request) {
	req := dto.GetUserByIDRequest{}
	if err := httpx.Parse(r, &req); err != nil {
		httpx.HttpResponseErr(w, err)
		return
	}
	res, err := c.UserSrv.GetUserID(r.Context(), &req)
	if err != nil {
		httpx.HttpResponseErr(w, err)
		return
	}
	httpx.HttpResponse(w, 200, res)
}

func (c *userControllerImpl) RawListUser(w http.ResponseWriter, r *http.Request) {
	req := dto.ListUserRequest{}
	if err := httpx.Parse(r, &req); err != nil {
		httpx.HttpResponseErr(w, err)
		return
	}
	res, err := c.UserSrv.ListUser(r.Context(), &req)
	if err != nil {
		httpx.HttpResponseErr(w, err)
		return
	}
	httpx.HttpResponse(w, 200, res)
}
