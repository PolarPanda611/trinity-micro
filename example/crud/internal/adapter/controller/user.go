// Author: Daniel TAN
// Date: 2021-08-18 23:39:51
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 01:27:48
// FilePath: /trinity-micro/example/crud/internal/adapter/controller/user.go
// Description:
package controller

import (
	"context"
	"net/http"
	"sync"

	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/core/dbx"
	"github.com/PolarPanda611/trinity-micro/core/logx"
	"github.com/PolarPanda611/trinity-micro/core/tracerx"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/dto"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/service"
	"github.com/go-chi/chi/v5"
)

func init() {
	UserControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(userControllerImpl)
		},
	}
	trinity.RegisterInstance("UserController", UserControllerPool)
	trinity.RegisterController("/example-api/v1/{corpID}/users", "UserController",
		trinity.NewRequestMapping("GET", "/", "ListUser",
			logx.ChiLoggerRequest,
			tracerx.ChiOpenTracer(
				tracerx.OperationNameFunc(
					func(r *http.Request) string {
						chiCtx := chi.RouteContext(r.Context())
						return r.Method + "=>" + chiCtx.RoutePattern()
					},
				),
			),
			dbx.SessionHandler,
		),
		trinity.NewRequestMapping("GET", "/{id}", "GetUserByID",
			logx.ChiLoggerRequest,
			tracerx.ChiOpenTracer(
				tracerx.OperationNameFunc(
					func(r *http.Request) string {
						chiCtx := chi.RouteContext(r.Context())
						return r.Method + "=>" + chiCtx.RoutePattern()
					},
				),
			),
			dbx.SessionHandler,
		),
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
